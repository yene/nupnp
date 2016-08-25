if (!window.location.origin) {
  window.location.origin = window.location.protocol + "//" + window.location.hostname + (window.location.port ? ':' + window.location.port: '');
}

document.addEventListener('DOMContentLoaded', setupUI);
function setupUI() {
  listDevices();
  return;

  var scanBtn = document.querySelector('.setup-scan');
  scanBtn.addEventListener('click', function(e) {
    e.preventDefault();
    scanBtn.classList.add('is-loading');
    console.log('start scaning');
    scanLights(function() {
      scanBtn.classList.remove('is-loading');
      listLights();
    });
  });
}

function listDevices() {
  var request = new XMLHttpRequest();
  request.open('GET', location.origin + '/devices.json', true);

  request.onload = function() {
    if (this.status >= 200 && this.status < 400) {
      var devices = JSON.parse(this.response);

      if (devices.length === 0) {
        return;
      }

      var list = document.querySelector('.device-list');
      list.innerHTML = '';

      devices.forEach(function(l) {

        var t = document.querySelector('.device-template');
        t.content.querySelector('.device-name').textContent = l.name;
        t.content.querySelector('.device-id').textContent = l.id;
        t.content.querySelector('.device-link').href = 'http://' + l.internaladdress;
        var clone = document.importNode(t.content, true);
        list.appendChild(clone);
      });

    } else {
      // We reached our target server, but it returned an error
    }
  };

  request.onerror = function() {
    // There was a connection error of some sort
  };

  request.send();
}

function addTD(e, text) {
  var td = document.createElement('td');
  td.innerHTML = text;
  return e.appendChild(td);
}

function addIcon(e, icon) {
  var i = document.createElement('i');
  i.classList.add('fa');
  i.classList.add('fa-' + icon);
  return e.appendChild(i);
}

function turnon(uuid) {
  var light = document.querySelector('.light-container[data-uuid="'+uuid+'"]');
  light.querySelector('.light-on').classList.add('is-loading');

  setTimeout(function(){
    light.querySelector('.light-on').classList.remove('is-loading');
    light.querySelector('.light-icon').dataset['state'] = String(true);
  }, 3000);

}

function turnoff(uuid, cb) {
  var light = document.querySelector('.light-container[data-uuid="'+uuid+'"]');
  light.querySelector('.light-off').classList.add('is-loading');

  setTimeout(function(){
    light.querySelector('.light-off').classList.remove('is-loading');
    light.querySelector('.light-icon').dataset['state'] = String(false);
  }, 3000);
}
