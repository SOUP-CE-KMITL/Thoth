var watson = require('watson-developer-cloud');
var fs = require('fs');
var speech_to_text = watson.speech_to_text({
  username: '0648b905-6758-4ae0-9c15-5cc53fadfa24',
  password: '9e90DAlqVwoK',
  version: 'v1'
});

var params = {
  audio: fs.createReadStream('sound.wav'),
  content_type: 'audio/wav'
};

speech_to_text.recognize(params, function(err, transcript) {
  if (err)
    console.log(err);
  else
    console.log(JSON.stringify(transcript, null, 2));
});