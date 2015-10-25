//var cp = Npm.require('child-process');
var Fiber = Npm.require('fibers');
var fs = Npm.require('fs');
var path = Npm.require('path');

dg = {};

dg.export = function(fileName, data, cb) {
   var hasError = false;
   var ok = true;
   var w =fs.createWriteStream(path.join('/tmp/dg',fileName));

   Fiber(function() {
      w.on('open',function() {
         data.forEach(function(v) {
            if(hasError) { return; }
            w.write(v.name+'\n');
         });

         w.end();
      });

      w.on('error', function(err) {
         hasError = true;
         return cb && cb(err);
      });

      w.on('finish',function() {
         if(hasError) { return; }
         return cb && cb();
      });
   });
};
