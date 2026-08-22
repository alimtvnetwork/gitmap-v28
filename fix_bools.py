import sys

replacements = [
    (r'gitmap/clonefrom/jsonschema_test.go', r'isNotNumber := !isNumber', r'isInvalidNumber := isNumber == false'),
    (r'gitmap/clonefrom/jsonschema_test.go', r'if isNotNumber == true', r'if isInvalidNumber'),
    (r'gitmap/cluster/distribution.go', r'hasNoClients := len\(clients\) == EmptySize', r'isEmpty := len(clients) == EmptySize'),
    (r'gitmap/cluster/distribution.go', r'if hasNoClients == true', r'if isEmpty'),
    (r'gitmap/cluster/exec_proj.go', r'isNotFound := foundPath == emptyString', r'isMissing := foundPath == emptyString'),
    (r'gitmap/cluster/exec_proj.go', r'if isNotFound == true', r'if isMissing'),
    (r'gitmap/cluster/node_resolver.go', r'hasNoNodes := len\(allNodes\) == constants.EmptySliceLength', r'isEmpty := len(allNodes) == constants.EmptySliceLength'),
    (r'gitmap/cluster/node_resolver.go', r'if hasNoNodes == true', r'if isEmpty'),
    (r'gitmap/cluster/ui.go', r'hasNoClients := numClients == DefaultProgress', r'isEmpty := numClients == DefaultProgress'),
    (r'gitmap/cluster/ui.go', r'if hasNoClients == true', r'if isEmpty'),
    (r'gitmap/cmd/auditlegacy.go', r'isNotAuditScannable := !isAuditScannable\(path\)', r'isUnscannable := isAuditScannable(path) == false'),
    (r'gitmap/cmd/auditlegacy.go', r'if isNotAuditScannable == true', r'if isUnscannable'),
    (r'gitmap/cmd/cdops.go', r'hasNoGroup := len\(group\) == 0', r'isEmpty := len(group) == 0'),
    (r'gitmap/cmd/cdops.go', r'if hasNoGroup == true', r'if isEmpty'),
    (r'gitmap/cmd/chromeprofile.go', r'isNotRunning := !isRunning', r'isStopped := isRunning == false'),
    (r'gitmap/cmd/chromeprofile.go', r'if isNotRunning == true', r'if isStopped'),
    (r'gitmap/cmd/chromeprofile_merge.go', r'isNotKnownMergeWhat := !isKnownMergeWhat\(\*what\)', r'isUnknownMergeWhat := isKnownMergeWhat(*what) == false'),
    (r'gitmap/cmd/chromeprofile_merge.go', r'if isNotKnownMergeWhat == true', r'if isUnknownMergeWhat'),
    (r'gitmap/cmd/chrome_bookmarks.go', r'hasNoRootsAfterRootName := rootName != "" && len\(roots\) == 0', r'isMissingRoots := rootName != "" && len(roots) == 0'),
    (r'gitmap/cmd/chrome_bookmarks.go', r'if hasNoRootsAfterRootName == true', r'if isMissingRoots'),
    (r'gitmap/cmd/chrome_bookmarks.go', r'hasNoRootsAfterFolderPath := folderPath != "" && len\(roots\) == 0', r'isMissingFolderRoots := folderPath != "" && len(roots) == 0'),
    (r'gitmap/cmd/chrome_bookmarks.go', r'if hasNoRootsAfterFolderPath == true', r'if isMissingFolderRoots'),
    (r'gitmap/cmd/chrome_bookmarks.go', r'hasNoRootsAfterMatch := hasMatchOrTitle == true && len\(roots\) == 0', r'isMissingMatchRoots := hasMatchOrTitle && len(roots) == 0'),
    (r'gitmap/cmd/chrome_bookmarks.go', r'if hasNoRootsAfterMatch == true', r'if isMissingMatchRoots'),
    (r'gitmap/cmd/clone.go', r'isNotDirectURL := !isDirectURL\(in\)', r'isIndirectURL := isDirectURL(in) == false'),
    (r'gitmap/cmd/clone.go', r'if isNotDirectURL == true', r'if isIndirectURL'),
]

import re
import os

for filepath, old, new in replacements:
    path = filepath.strip()
    if not os.path.exists(path):
        continue
    with open(path, 'r', encoding='utf-8') as f:
        content = f.read()
    content = re.sub(old, new, content)
    with open(path, 'w', encoding='utf-8') as f:
        f.write(content)

print("Applied 13 replacements!")
