const fetch = require('node-fetch');
const tar = require('tar');
const zlib = require('zlib');

const install = async (callback) => {
    let gunzip = zlib.createGunzip();
    let untar = tar.x({ cwd: '../../.bin' });

    // Pipe error to callback
    gunzip.on('error', callback);
    untar.on('error', callback);

    let req = await fetch(
        'https://jo-compiler.s3.eu-central-1.amazonaws.com/elljo_0.0.0-SNAPSHOT-acf1211_darwin_amd64.tar.gz'
    );
    req.body.pipe(gunzip).pipe(untar);
};

function uninstall(callback) {}

let argv = process.argv;
if (argv && argv.length > 2) {
    let cmd = process.argv[2];
    if (cmd === 'install') {
        install((err) => {
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
