package cmdcdp

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// WSClient provides a lightweight RFC 6455 WebSocket client using pure Go standard library.
type WSClient struct {
	conn   net.Conn
	reader *bufio.Reader
}

func resolveHostPort(u *url.URL) string {
	if strings.Contains(u.Host, ":") {
		return u.Host
	}

	if u.Scheme == "wss" {
		return u.Host + ":443"
	}

	return u.Host + ":80"
}

func dialConnection(u *url.URL, timeout time.Duration) (net.Conn, error) {
	addr := resolveHostPort(u)
	if u.Scheme == "wss" {
		dialer := &net.Dialer{Timeout: timeout}

		return tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{MinVersion: tls.VersionTLS12})
	}

	return net.DialTimeout("tcp", addr, timeout)
}

func generateWebSocketKey() (string, error) {
	keyBytes := make([]byte, 16)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(keyBytes), nil
}

func buildHandshakeRequest(u *url.URL, secKey string) string {
	path := u.RequestURI()
	if path == "" {
		path = "/"
	}

	return fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: %s\r\nSec-WebSocket-Version: 13\r\n\r\n", path, u.Host, secKey)
}

func discardResponseHeaders(reader *bufio.Reader) error {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return apperror.WrapSimple(err, "ws.handshake.headers")
		}

		if strings.TrimSpace(line) == "" {
			return nil
		}
	}
}

func verifyHandshake(reader *bufio.Reader) error {
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		return apperror.WrapSimple(err, "ws.handshake.status")
	}

	if !strings.Contains(statusLine, "101") {
		return fmt.Errorf("unexpected websocket status: %s", strings.TrimSpace(statusLine))
	}

	return discardResponseHeaders(reader)
}

func sendHandshakeRequest(conn net.Conn, u *url.URL) error {
	secKey, err := generateWebSocketKey()
	if err != nil {
		return apperror.WrapSimple(err, "ws.key")
	}

	req := buildHandshakeRequest(u, secKey)
	_, writeErr := conn.Write([]byte(req))
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "ws.write_req")
	}

	return nil
}

func performHandshake(conn net.Conn, u *url.URL) (*WSClient, error) {
	if err := sendHandshakeRequest(conn, u); err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	if err := verifyHandshake(reader); err != nil {
		return nil, err
	}

	return &WSClient{conn: conn, reader: reader}, nil
}

func connectWebSocket(u *url.URL, timeout time.Duration) (*WSClient, error) {
	conn, err := dialConnection(u, timeout)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ws.dial")
	}

	client, err := performHandshake(conn, u)
	if err != nil {
		conn.Close()

		return nil, err
	}

	return client, nil
}

// DialWebSocket connects to a WebSocket server and completes the RFC 6455 upgrade handshake.
func DialWebSocket(urlStr string, timeout time.Duration) (*WSClient, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ws.parse_url")
	}

	return connectWebSocket(u, timeout)
}

func buildExtended64Header(payloadLen int) []byte {
	hdr := make([]byte, 10)
	hdr[0] = 0x81
	hdr[1] = 0x80 | 127
	binary.BigEndian.PutUint64(hdr[2:10], uint64(payloadLen))

	return hdr
}

func buildUint16Header(payloadLen int) []byte {
	hdr := make([]byte, 4)
	hdr[0] = 0x81
	hdr[1] = 0x80 | 126
	binary.BigEndian.PutUint16(hdr[2:4], uint16(payloadLen))

	return hdr
}

func buildFrameHeader(payloadLen int) []byte {
	if payloadLen < 126 {
		return []byte{0x81, 0x80 | byte(payloadLen)}
	}

	if payloadLen <= 65535 {
		return buildUint16Header(payloadLen)
	}

	return buildExtended64Header(payloadLen)
}

func generateMaskKey() ([4]byte, error) {
	var mask [4]byte
	_, err := rand.Read(mask[:])
	if err != nil {
		return mask, err
	}

	return mask, nil
}

func maskPayload(payload []byte, mask [4]byte) []byte {
	masked := make([]byte, len(payload))
	for i := range payload {
		masked[i] = payload[i] ^ mask[i%4]
	}

	return masked
}

func constructMaskedFrame(payload []byte, mask [4]byte) []byte {
	hdr := buildFrameHeader(len(payload))
	masked := maskPayload(payload, mask)
	frame := append(hdr, mask[:]...)

	return append(frame, masked...)
}

// SendFrame transmits a text frame masked with a 4-byte XOR mask per RFC 6455.
func (c *WSClient) SendFrame(payload []byte) error {
	mask, err := generateMaskKey()
	if err != nil {
		return apperror.WrapSimple(err, "ws.mask_gen")
	}

	frame := constructMaskedFrame(payload, mask)
	_, writeErr := c.conn.Write(frame)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "ws.write_frame")
	}

	return nil
}

func readUint16Length(reader *bufio.Reader) (int, error) {
	var buf [2]byte
	if _, err := io.ReadFull(reader, buf[:]); err != nil {
		return 0, apperror.WrapSimple(err, "ws.read_u16")
	}

	return int(binary.BigEndian.Uint16(buf[:])), nil
}

func readUint64Length(reader *bufio.Reader) (int, error) {
	var buf [8]byte
	if _, err := io.ReadFull(reader, buf[:]); err != nil {
		return 0, apperror.WrapSimple(err, "ws.read_u64")
	}

	return int(binary.BigEndian.Uint64(buf[:])), nil
}

func readPayloadLength(reader *bufio.Reader, lenField int) (int, error) {
	if lenField < 126 {
		return lenField, nil
	}

	if lenField == 126 {
		return readUint16Length(reader)
	}

	return readUint64Length(reader)
}

func unmaskIfRequired(payload []byte, mask [4]byte, hasMask bool) []byte {
	if !hasMask {
		return payload
	}

	return maskPayload(payload, mask)
}

func readMaskIfRequired(reader *bufio.Reader, hasMask bool) ([4]byte, error) {
	var mask [4]byte
	if !hasMask {
		return mask, nil
	}

	_, err := io.ReadFull(reader, mask[:])
	if err != nil {
		return mask, apperror.WrapSimple(err, "ws.read_mask")
	}

	return mask, nil
}

func (c *WSClient) readPayloadBytes(length int, hasMask bool) ([]byte, error) {
	mask, err := readMaskIfRequired(c.reader, hasMask)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(c.reader, payload); err != nil {
		return nil, apperror.WrapSimple(err, "ws.read_payload")
	}

	return unmaskIfRequired(payload, mask, hasMask), nil
}

func parseFrameHeader(reader *bufio.Reader) (byte, bool, int, error) {
	var header [2]byte
	if _, err := io.ReadFull(reader, header[:]); err != nil {
		return 0, false, 0, apperror.WrapSimple(err, "ws.read_hdr")
	}

	opcode := header[0] & 0x0F
	hasMask := (header[1] & 0x80) != 0
	length, err := readPayloadLength(reader, int(header[1]&0x7F))

	return opcode, hasMask, length, err
}

func (c *WSClient) readSingleFrame() ([]byte, byte, error) {
	opcode, hasMask, length, err := parseFrameHeader(c.reader)
	if err != nil {
		return nil, 0, err
	}

	payload, err := c.readPayloadBytes(length, hasMask)

	return payload, opcode, err
}

func (c *WSClient) handlePing(payload []byte) error {
	mask, err := generateMaskKey()
	if err != nil {
		return err
	}

	hdr := []byte{0x8A, 0x80 | byte(len(payload)&0x7F)}
	masked := maskPayload(payload, mask)
	frame := append(append(hdr, mask[:]...), masked...)

	_, writeErr := c.conn.Write(frame)

	return writeErr
}

func (c *WSClient) processReceivedFrame(payload []byte, opcode byte) ([]byte, bool, error) {
	if opcode == 0x8 {
		return nil, false, io.EOF
	}

	if opcode == 0x9 {
		return nil, false, c.handlePing(payload)
	}

	return payload, true, nil
}

func (c *WSClient) readNextValidPayload() ([]byte, bool, error) {
	payload, opcode, err := c.readSingleFrame()
	if err != nil {
		return nil, false, err
	}

	return c.processReceivedFrame(payload, opcode)
}

// RecvFrame reads and decodes the next WebSocket frame, resolving extended lengths and unmasking.
func (c *WSClient) RecvFrame() ([]byte, error) {
	for {
		res, isDone, err := c.readNextValidPayload()
		if err != nil {
			return nil, err
		}

		if isDone {
			return res, nil
		}
	}
}

// Close closes the underlying network connection.
func (c *WSClient) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
