import fs from 'fs';
import path from 'path';

const TARGET_EXTENSIONS = ['.ts', '.tsx', '.php', '.py'];
const IGNORE_DIRS = ['node_modules', '.git'];

let hasViolations = false;

function scanDir(dir) {
    if (!fs.existsSync(dir)) return;
    const files = fs.readdirSync(dir);
    for (const file of files) {
        if (IGNORE_DIRS.includes(file)) continue;
        
        const fullPath = path.join(dir, file);
        let stat;
        try {
            stat = fs.statSync(fullPath);
        } catch (e) {
            continue;
        }
        
        if (stat.isDirectory()) {
            scanDir(fullPath);
        } else if (stat.isFile()) {
            const ext = path.extname(file);
            if (TARGET_EXTENSIONS.includes(ext)) {
                checkFile(fullPath, ext);
            }
        }
    }
}

function checkFile(filePath, ext) {
    const content = fs.readFileSync(filePath, 'utf-8');
    const lines = content.split(/\r?\n/);
    
    lines.forEach((line, index) => {
        const lineNumber = index + 1;
        
        // 1. No TypeScript string unions for status/state
        if (ext === '.ts' || ext === '.tsx') {
            if (/((?:type\s+\w*(?:Status|State)\s*=)|(?:status|state)\s*\??:\s*)(['"][^'"]+['"]\s*\|\s*['"][^'"]+['"])/i.test(line)) {
                reportViolation(filePath, lineNumber, 'No TypeScript string unions for status/state');
            } else if (/type\s+\w+\s*=\s*(['"][^'"]+['"]\s*\|\s*['"][^'"]+['"])/.test(line) && /status|state/i.test(line)) {
                 reportViolation(filePath, lineNumber, 'No TypeScript string unions for status/state');
            }
        }
        
        // 2. All Enums must end in Type
        if (ext === '.ts' || ext === '.tsx' || ext === '.php') {
            if (/^\s*(?:export\s+)?enum\s+(?!\w+Type\b)\w+\b/.test(line)) {
                reportViolation(filePath, lineNumber, 'Enums must have names ending in `Type`');
            }
        }
        
        // 3. No inverted success checks
        if (/!\s*[a-zA-Z0-9_$->]*\bisSuccess\b/.test(line)) {
            reportViolation(filePath, lineNumber, 'No inverted success checks (e.g., !isSuccess, !response.isSuccess())');
        }
        
        // 4. (PHP) No $isFailed = !$exists; assignments
        if (ext === '.php') {
            if (/\$isFailed\s*=\s*!\s*\$exists\b/.test(line)) {
                reportViolation(filePath, lineNumber, 'No `$isFailed = !$exists;` assignments');
            }
        }
    });
}

function reportViolation(filePath, lineNumber, message) {
    console.error(`Violation in ${filePath}:${lineNumber} - ${message}`);
    hasViolations = true;
}

const rootDirs = ['src', 'spec'];
const cwd = process.cwd();

let scannedSomething = false;
for (const dir of rootDirs) {
    const fullDir = path.join(cwd, dir);
    if (fs.existsSync(fullDir)) {
        scanDir(fullDir);
        scannedSomething = true;
    }
}

// Fallback to scanning the current directory if src/ and spec/ are not found.
if (!scannedSomething) {
    scanDir(cwd);
}

if (hasViolations) {
    process.exit(1);
} else {
    process.exit(0);
}
