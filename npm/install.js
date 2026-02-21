const https = require('https');
const fs = require('fs');
const path = require('path');

const packageJson = require('./package.json');
const version = packageJson.version;

const platform = process.platform;
const arch = process.arch;

// 플랫폼별 파일명
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

const downloadUrl = `https://github.com/indigo-net/Brf.it/releases/download/v${version}/${binaryName}`;
const binDir = path.join(__dirname, 'bin');
const outputPath = path.join(binDir, binaryName);

// bin 디렉토리 생성
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

console.log(`Downloading brfit v${version} for ${platform}-${arch}...`);

const file = fs.createWriteStream(outputPath);

https.get(downloadUrl, (response) => {
  if (response.statusCode === 302 || response.statusCode === 301) {
    // 리다이렉트 처리
    https.get(response.headers.location, (redirectResponse) => {
      redirectResponse.pipe(file);
    }).on('error', (err) => {
      console.error('Download failed:', err.message);
      process.exit(1);
    });
  } else if (response.statusCode === 200) {
    response.pipe(file);
  } else {
    console.error(`Download failed: HTTP ${response.statusCode}`);
    console.error('Make sure the release exists at:', downloadUrl);
    process.exit(1);
  }
}).on('error', (err) => {
  console.error('Download failed:', err.message);
  process.exit(1);
});

file.on('finish', () => {
  file.close();
  // 실행 권한 부여 (Unix)
  if (platform !== 'win32') {
    fs.chmodSync(outputPath, '755');
  }
  console.log('brfit installed successfully!');
});
