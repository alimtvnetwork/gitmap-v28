import fs from 'fs';
import path from 'path';
import ts from 'typescript';

const IGNORE_DIRS = new Set([
  'node_modules', 'dist', 'build', '.git', '.next', 'coverage',
  '.gemini', '.system_generated', 'temp', 'scratch'
]);

function getFiles(dir) {
  let results = [];
  if (!fs.existsSync(dir)) return results;
  const list = fs.readdirSync(dir);
  for (const file of list) {
    if (IGNORE_DIRS.has(file)) continue;
    const fullPath = path.join(dir, file);
    try {
      const stat = fs.statSync(fullPath);
      if (stat.isDirectory()) {
        results = results.concat(getFiles(fullPath));
      } else if (file.endsWith('.ts') || file.endsWith('.tsx')) {
        results.push(fullPath);
      }
    } catch {
      // ignore unreadable files
    }
  }
  return results;
}

const targetDirs = ['src'];
const files = targetDirs.flatMap(d => getFiles(d));
let violations = [];

files.forEach(file => {
  const content = fs.readFileSync(file, 'utf8');
  const sf = ts.createSourceFile(file, content, ts.ScriptTarget.Latest, true);

  function checkNode(node, inIf) {
    if (ts.isIfStatement(node)) {
      if (inIf) {
        const { line } = sf.getLineAndCharacterOfPosition(node.getStart());
        violations.push({
          file: path.relative('.', file).replace(/\\/g, '/'),
          line: line + 1,
          text: node.getText(sf).split('\n')[0].trim()
        });
      }
      ts.forEachChild(node.thenStatement, child => {
        if (ts.isFunctionDeclaration(child) || ts.isArrowFunction(child) || ts.isFunctionExpression(child)) {
          checkNode(child, false);
        } else {
          checkNode(child, true);
        }
      });
      if (node.elseStatement) {
        if (ts.isIfStatement(node.elseStatement)) {
          checkNode(node.elseStatement, inIf);
        } else {
          ts.forEachChild(node.elseStatement, child => {
            if (ts.isFunctionDeclaration(child) || ts.isArrowFunction(child) || ts.isFunctionExpression(child)) {
              checkNode(child, false);
            } else {
              checkNode(child, true);
            }
          });
        }
      }
      return;
    }

    if (ts.isFunctionDeclaration(node) || ts.isArrowFunction(node) || ts.isFunctionExpression(node)) {
      ts.forEachChild(node, child => checkNode(child, false));
      return;
    }

    ts.forEachChild(node, child => checkNode(child, inIf));
  }

  checkNode(sf, false);
});

if (violations.length > 0) {
  console.error(`❌ FAIL: ${violations.length} nested if violation(s) found in TypeScript files:`);
  violations.forEach(v => {
    console.error(`  ${v.file}:${v.line}: Nested if statement found: ${v.text}`);
  });
  process.exit(1);
} else {
  console.log(`✅ PASS: Zero nested if statements found across ${files.length} TypeScript file(s).`);
  process.exit(0);
}
