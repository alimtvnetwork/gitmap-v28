import os

def replace_in_file(path, replacements, add_import=True):
    if not os.path.exists(path): return
    with open(path, "r", encoding="utf-8") as f:
        c = f.read()
    
    for old, new in replacements:
        c = c.replace(old, new)
        
    if add_import and '"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"' not in c:
        c = c.replace('import (', 'import (\n\t"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"\n', 1)
        
    with open(path, "w", encoding="utf-8") as f:
        f.write(c)

# create.go
replace_in_file("gitmap/archive/create.go", [
    ('fmt.Errorf("%w: %q", ErrUnknownFormat, path)', 'apperror.Wrap(ErrUnknownFormat, "validate format", map[string]any{"path": path})'),
    ('fmt.Errorf("%s archives are read-only in this build (use zip or tar.*)", format)', 'apperror.New("validate format", "ERR_READONLY_FORMAT", map[string]any{"format": format})'),
    ('fmt.Errorf("gather sources: %w", err)', 'apperror.Wrap(err, "gather sources", nil)')
])

# extract.go
replace_in_file("gitmap/archive/extract.go", [
    ('fmt.Errorf("identify %q: %w", srcArchive, err)', 'apperror.Wrap(err, "identify archive", map[string]any{"src": srcArchive})'),
    ('fmt.Errorf("extract: %w", err)', 'apperror.Wrap(err, "extract", nil)'),
    ('fmt.Errorf("identify: %w", err)', 'apperror.Wrap(err, "identify", nil)'),
    ('fmt.Errorf("format %s is not extractable", format.Extension())', 'apperror.New("extract", "ERR_UNSUPPORTED_FORMAT", map[string]any{"format": format.Extension()})'),
    ('fmt.Errorf("rejecting entry with unsafe path: %q", entry.NameInArchive)', 'apperror.New("extract entry", "ERR_UNSAFE_PATH", map[string]any{"path": entry.NameInArchive})')
])

# list.go
replace_in_file("gitmap/archive/list.go", [
    ('fmt.Errorf("archive identify: %w", err)', 'apperror.Wrap(err, "identify archive", nil)'),
    ('fmt.Errorf("format %s is not extractable", format.Extension())', 'apperror.New("list archive", "ERR_UNSUPPORTED_FORMAT", map[string]any{"format": format.Extension()})')
])

# source.go
replace_in_file("gitmap/archive/source.go", [
    ('fmt.Errorf("local source %q: %w", raw, err)', 'apperror.Wrap(err, "resolve local source", map[string]any{"raw": raw})'),
    ('fmt.Errorf("download %q: %w", raw, err)', 'apperror.Wrap(err, "download source", map[string]any{"raw": raw})'),
    ('fmt.Errorf("http %s", resp.Status)', 'apperror.New("download source", "ERR_HTTP", map[string]any{"status": resp.Status})'),
    ('fmt.Errorf("git clone %q: %w", raw, err)', 'apperror.Wrap(err, "git clone", map[string]any{"raw": raw})'),
    ('fmt.Errorf("no archive in %s", abs)', 'apperror.New("find archive", "ERR_NOT_FOUND", map[string]any{"dir": abs})'),
    ('fmt.Errorf("found %d archives in %s", len(found), abs)', 'apperror.New("find archive", "ERR_MULTIPLE_FOUND", map[string]any{"count": len(found), "dir": abs})')
])

print("Replaced fmt.Errorf with apperror")
