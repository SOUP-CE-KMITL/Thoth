var speech=null;
var audio_context,
    recorder,
    volume,
    volumeLevel = 0,
    currentEditedSoundIndex;

function startUserMedia(stream) {
  var input = audio_context.createMediaStreamSource(stream);
  console.log('Media stream created.');

  volume = audio_context.createGain();
  volume.gain.value = volumeLevel;
  input.connect(volume);
  volume.connect(audio_context.destination);
  console.log('Input connected to audio context destination.');
  
  recorder = new Recorder(input);
  console.log('Recorder initialised.');
}

function changeVolume(value) {
  if (!volume) return;
  volumeLevel = value;
  volume.gain.value = value;
}

function startRecording(button) {
  recorder && recorder.record();
  button.disabled = true;
  button.nextElementSibling.disabled = false;
  console.log('Recording...');
}

function stopRecording(button) {
  recorder && recorder.stop();
  button.disabled = true;
  button.previousElementSibling.disabled = false;
  console.log('Stopped recording.');
  
  // create WAV download link using audio data blob
  createDownloadLink();
  
  recorder.clear();
}

function createDownloadLink() {
  currentEditedSoundIndex = -1;
  recorder && recorder.exportWAV(handleWAV.bind(this));
}

function handleWAV(blob) {
  var tableRef = document.getElementById('recordingslist');
  if (currentEditedSoundIndex !== -1) {
    $('#recordingslist tr:nth-child(' + (currentEditedSoundIndex + 1) + ')').remove();
  }
  speech=blob;
  
  PostToWatson();
  /*
  var url = URL.createObjectURL(blob);
  var newRow   = tableRef.insertRow(currentEditedSoundIndex);
  newRow.className = 'soundBite';
  var audioElement = document.createElement('audio');
  var downloadAnchor = document.createElement('a');
  var editButton = document.createElement('button');
  
  audioElement.controls = true;
  audioElement.src = url;

  downloadAnchor.href = url;
  downloadAnchor.download = new Date().toISOString() + '.wav';
  downloadAnchor.innerHTML = 'Download';
  downloadAnchor.className = 'btn btn-primary';
  var newCell = newRow.insertCell(-1);
  newCell.appendChild(audioElement);
  newCell = newRow.insertCell(-1);
  newCell.appendChild(downloadAnchor);
  newCell = newRow.insertCell(-1);
*/
}

function PostToWatson() {
	$.ajax({
		url: "https://paas.jigko.net/speech/",
	type: "POST",
	contentType: "audio/wav",
	data: speech,
	processData: false
	}).done(function(result) {
		console.log(result);
		/*
		   {
		   "results": [
		   	{
		   		"alternatives": [
		   			{
				   	"confidence": 0.994, 
				   	"transcript": "hello "
				   	}
		   		], 
			   	"final": true
		   	}
		   ], 
		   "result_index": 0
		   }
		*/
		resJson = JSON.parse(result);
		text = resJson.results[0].alternatives[0].transcript;
		console.log(text);
		
		// ------------Command-------------
		// Next,Back,App,Node,error(era)
		if (text.indexOf("app")>=0){
			//alert("Apps");
			window.location.href = "https://paas.jigko.net/node/"
		}else if (text.indexOf("node")>=0||text.indexOf("nose")>=0||text.indexOf("north")>=0||text.indexOf("no")>=0){
			//alert("Nodes");
			window.location.href = "https://paas.jigko.net/nodes/"
		}else if (text.indexOf("error")>=0||text.indexOf("era")>=0){
			alert("Error Apps");
		}else if (text.indexOf("next")>=0){
			window.history.next();
		}else if (text.indexOf("back")>=0){
			window.history.back();
		}else{
			alert("I can't understand your command");
		}
	});
}

window.onload = function init() {
  try {
    // webkit shim
    window.AudioContext = window.AudioContext || window.webkitAudioContext || window.mozAudioContext;
    navigator.getUserMedia = navigator.getUserMedia || navigator.webkitGetUserMedia || navigator.mozGetUserMedia;
    window.URL = window.URL || window.webkitURL || window.mozURL;
    
    audio_context = new AudioContext();
    console.log('Audio context set up.');
    console.log('navigator.getUserMedia ' + (navigator.getUserMedia ? 'available.' : 'not present!'));
  } catch (e) {
    console.warn('No web audio support in this browser!');
  }
  
  navigator.getUserMedia({audio: true}, startUserMedia, function(e) {
    console.warn('No live audio input: ' + e);
  });
};
