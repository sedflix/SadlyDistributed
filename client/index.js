// knockout
var vm = {
    runningATask: ko.observable(false),
    tasksDone: ko.observable(0)
}
ko.applyBindings(vm)

// other stuff
function noop() {}


memoized_wasm = {}
function getCompiledWasmPromise(path) {
  if (!memoized_wasm[path]) {
    return WebAssembly.compileStreaming(fetch(path)).then((module) => memoized_wasm[path] = module)
  } else {
    return Promise.resolve(memoized_wasm[path])
  }
}

var s = gotalk.connection().on('open', function () {
    console.log('socket opened');
    s.request("resource-available", {cores: navigator.hardwareConcurrency, mem: 512}, function (err, result) {
        if (err) return console.log('request for "resource-available" failed:', err);
        console.log('', result);
    });


    gotalk.handle("receive-job", function (params, result) {

      console.log("got a job: ", params);
      vm.runningATask(true);

      setTimeout(function () {
        instantiate_go_script(params)
      }, 0);

      result({"is_okay": "Okay"});
    });
}).on('close', function (err) {
    console.log('socket closed', err);
}).on('notification', function (n, v) {
    console.log('received notification "' + name + '":', value);
});

function hijackOutput(cb) {
  var log = console.log;
  var output = "";
  console.log = function () {
    var args = Array.prototype.slice.call(arguments);
    output += args.map(function (x) {
        return String(x)
    });
    log.apply(this, args);
  };
  cb(log);
  console.log = log;
  return output;
}

function instantiate_go_script(params) {
  const go = new Go();
  getCompiledWasmPromise(params.wasm)
    .then(function (compiled) {
      //console.log("got compiled", compiled)
      return WebAssembly.instantiate(compiled, go.importObject)
    })
    .then((instantiated) => {
      //console.log("got instantiated", instantiated)
      var output = hijackOutput(function () {

        go.argv = ["js"].concat(params.parameters.split(" "));
        go.run(instantiated);
      });

      s.request("job-complete", {
        job_id: params.job_id,
        program_id: params.program_id,
        parameters: params.parameters,
        result: output
      }, noop);

      console.log("finished job", output);
      vm.runningATask(false);
      vm.tasksDone(vm.tasksDone() + 1)
    }
  )
  .catch(x => console.log("error", x));
};

function instantiate_emscripten(params) {
  var Module = {
    preRun: [],
    postRun: [],
    print: (function() {
      return function(text) {
        if (arguments.length > 1) text = Array.prototype.slice.call(arguments).join(' ');
        console.log(text);
      };
    })(),
    printErr: function(text) {
      if (arguments.length > 1) text = Array.prototype.slice.call(arguments).join(' ');
      if (0) { // XXX disabled for safety typeof dump == 'function') {
        dump(text + '\n'); // fast, straight to the real console
      } else {
        console.error(text);
      }
    },
    setStatus: function(text) {
      if (!Module.setStatus.last) Module.setStatus.last = { time: Date.now(), text: '' };
      if (text === Module.setStatus.last.text) return;
      var m = text.match(/([^(]+)\((\d+(\.\d+)?)\/(\d+)\)/);
      var now = Date.now();
      if (m && now - Module.setStatus.last.time < 30) return; // if this is a progress update, skip it if too soon
      Module.setStatus.last.time = now;
      Module.setStatus.last.text = text;
      if (m) {
        text = m[1];
      }
    },
    totalDependencies: 0,
    monitorRunDependencies: function(left) {
      this.totalDependencies = Math.max(this.totalDependencies, left);
      Module.setStatus(left ? 'Preparing... (' + (this.totalDependencies-left) + '/' + this.totalDependencies + ')' : 'All downloads complete.');
    }
  };
  Module.setStatus('Downloading...');
  window.onerror = function(event) {
    // TODO: do not warn on ok events like simulating an infinite loop or exitStatus
    Module.setStatus('Exception thrown, see JavaScript console');
    Module.setStatus = function(text) {
      if (text) Module.printErr('[post-exception status] ' + text);
    };
  };
}
