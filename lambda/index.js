var child_process = require('child_process');

exports.handler = function(event, context) {
  console.log(event)

  var proc = child_process.spawn('./collector', [], { stdio: [process.stdin, 'pipe', 'pipe'] });

  proc.stdout.on('data', function(line){
    var msg = JSON.parse(line);
    console.log("stdout: " + msg)
    context.succeed(msg);
  })

  proc.stderr.on('data', function(line){
    var msg = new Error(line)
    console.log("stderr: " + msg)
    context.fail(msg);
  })

  proc.on('exit', function(code){
    if (code == 0) {
      console.log("Ran successfully, but no action needed")
      context.succeed("Ran successfully, but no action needed")
    } else {
      console.error('exit: %s', code)
      context.fail("Ran unsuccessfully")
    }
  })
}
