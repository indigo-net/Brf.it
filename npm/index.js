#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');

const platform = process.platform;
const arch = process.arch;

// 플랫폼별 바이너리 경로 결정
let binaryName;
if (platform === 'win32') {
  binaryName = 'brfit-windows-amd64.exe';
} else if (platform === 'darwin') {
  binaryName = arch === 'arm64' ? 'brfit-darwin-arm64' : 'brfit-darwin-amd64';
} else if (platform === 'linux') {
  binaryName = arch === 'arm64' ? 'brfit-linux-arm64' : 'brfit-linux-amd64';
} else {
  console.error(`Unsupported platform: ${platform}-${arch}`);
  process.exit(1);
}

const binaryPath = path.join(__dirname, 'bin', binaryName);

// 바이너리 실행
const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit'
});

child.on('exit', (code) => {
  process.exit(code || 0);
});
