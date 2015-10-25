Package.describe({
   name: 'dg',
   version: '0.0.1',
   summary: 'Simple danguardian control interface',
   documentation: null
});

Npm.depends({
});

Package.onUse(function(api) {
   api.addFiles('dg.js', ['server']);
   api.export('dg',['server']);
});
