var xhr = new XMLHttpRequest();
datamap = new Map();
datamap.set('http_request', 'hello :)')
datamap.set('http_request2', 2)
request = {
  Cache: datamap,
  PrivateKey: 'hello',
}

xhr.open("POST", 'http://localhost/tkv_v1/save', true);
xhr.setRequestHeader('Content-Type', 'application/json');
xhr.send(JSON.stringify(request));