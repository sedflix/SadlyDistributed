function noop() {}

const go = new Go();

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


function instantiate_go_script(params) {
  WebAssembly.instantiateStreaming(fetch(params.wasm), go.importObject).then((result) => {
    setTimeout(function () {
      var log = console.log;
      var output = "";
      console.log = function () {
        var args = Array.prototype.slice.call(arguments);
        output += args.map(function (x) {
            return String(x)
        });
        log.apply(this, args);
      };

      go.argv = ["js"].concat(params.parameters.split(" "));
      go.run(result.instance);

      s.request("job-complete", {
        job_id: params.job_id,
        program_id: params.program_id,
        parameters: params.parameters,
        result: output
      }, noop);

      log("finished job", output);
      vm.runningATask(false);
      vm.tasksDone(vm.tasksDone() + 1)
    }, 0);
  });
};


