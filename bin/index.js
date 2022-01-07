const fetch = require('node-fetch');
const tar = require('tar');
const zlib = require('zlib');

const install = async callback => {
  let gunzip = zlib.createGunzip();
  let untar = tar.x({ cwd: '../../.bin' });

  // Pipe error to callback
  gunzip.on('error', callback);
  untar.on('error', callback);

  const architecture = process.arch;
  const os = process.platform;

  let url = 'https://elljo.s3.eu-central-1.amazonaws.com/elljo_0.0.2-alpha_';

  switch (os) {
    case 'darwin':
      url += 'darwin_';
      break;
    case 'linux':
      url += 'linux_';
      break;
    case 'win32':
      url += 'windows_';
      break;
    default:
      console.error('Your system is not supported by EllJo!');
      process.exit(1);
  }

  switch (architecture) {
    case 'x64':
      url += 'amd64';
      break;
    case 'arm64':
      url += 'arm64';
      break;
    default:
      console.error('Your cpu architecture is not supported by EllJo!');
      process.exit(1);
  }

  url += '.tar.gz';

  let req = await fetch(url);
  req.body.pipe(gunzip).pipe(untar);
};

function uninstall(callback) {}

let argv = process.argv;
if (argv && argv.length > 2) {
  let cmd = process.argv[2];
  if (cmd === 'install') {
    install(err => {
      if (err) {
        console.error(err);
        process.exit(1);
      } else {
        process.exit(0);
      }
    });
  } else if (cmd === 'uninstall') {
    // TODO: Implement uninstall
    uninstall();
  }
}