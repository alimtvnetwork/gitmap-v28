import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, '..');

const IGNORE_DIRS = new Set([
  'node_modules', '.git', 'dist', 'build', 'tmp', '.ci-out', '.github', 'release-artifacts', 'reports'
]);

function getFiles(dir, exts, fileList = []) {
  let files;
  try {
    files = fs.readdirSync(dir);
  } catch (e) {
    return fileList;
  }
  
  for (const file of files) {
    if (IGNORE_DIRS.has(file)) continue;
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);
    if (stat.isDirectory()) {
      getFiles(filePath, exts, fileList);
    } else {
      if (exts.includes(path.extname(filePath))) {
        fileList.push(filePath);
      }
    }
  }
  return fileList;
}

const exts = ['.go', '.ts', '.tsx', '.js'];
const files = getFiles(rootDir, exts);
const shouldFix = process.argv.includes('--fix');

let hasError = false;

function reportError(file, line, msg) {
  const relPath = path.relative(rootDir, file);
  console.error(`${relPath}:${line} - ${msg}`);
  hasError = true;
}

for (const file of files) {
  let content;
  try {
    content = fs.readFileSync(file, 'utf8');
  } catch (e) {
    continue;
  }
  const lines = content.split(/\r?\n/);
  let modified = false;

  for (let i = 0; i < lines.length; i++) {
    const text = lines[i].trim();
    
    // 1. No double empty lines (\n\n\n)
    if (i > 0 && text === '' && lines[i-1].trim() === '') {
      if (!(i > 1 && lines[i-2].trim() === '')) {
        reportError(file, i + 1, "Double empty lines (Rule R13 extension)");
        if (shouldFix) { lines.splice(i, 1); i--; modified = true; continue; }
      }
    }

    // 2. No empty line at the start of a function/block
    if (i > 0 && text === '' && lines[i-1].trim().endsWith('{')) {
      reportError(file, i + 1, "Empty line at the start of a block (Rule R12)");
      if (shouldFix) { lines.splice(i, 1); i--; modified = true; continue; }
    }

    // 3. Blank line before return
    if (/^return(\s|;|$)/.test(text)) {
      let j = i - 1;
      while (j >= 0 && lines[j].trim().startsWith('//')) {
        j--;
      }
      if (j >= 0) {
        const prev = lines[j].trim();
        if (prev !== '' && !prev.endsWith('{') && !prev.endsWith(':')) {
          reportError(file, i + 1, "Blank line required before return (Rule R4)");
          if (shouldFix) { lines.splice(i, 0, ''); i++; modified = true; continue; }
        }
      }
    }

    // 4. Blank line after } if followed by more code
    if (text === '}') {
      if (i + 1 < lines.length) {
        const next = lines[i+1].trim();
        if (next !== '' && 
            !/^[}\)\];,\.]/.test(next) &&
            !/^<\//.test(next) &&
            !/^(else|catch|finally|while|case|default)\b/.test(next)) {
          reportError(file, i + 1, "Blank line required after } when followed by more code (Rule R5)");
          if (shouldFix) { lines.splice(i + 1, 0, ''); modified = true; continue; }
        }
      }
    }
  }

  if (modified && shouldFix) {
    fs.writeFileSync(file, lines.join('\n'));
  }
}

if (hasError) {
  process.exit(1);
} else {
  process.exit(0);
}
