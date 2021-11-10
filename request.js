var xhr = new XMLHttpRequest();
datamap = new Map();
datamap.set('http_request', 'hello :)')
datamap.set('http_request2', 2)
request = {
  Cache: datamap,
  PrivateKey: '5d41402abc4b2a76b9719d911017c592',
}

xhr.open("POST", 'http://localhost/tkv_v1/save', true);
xhr.setRequestHeader('Content-Type', 'application/json');
xhr.send(JSON.stringify(request));