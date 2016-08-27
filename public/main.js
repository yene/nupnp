
// for IE 9,10,11
if (!window.location.origin) {
  window.location.origin = window.location.protocol + "//" + window.location.hostname + (window.location.port ? ':' + window.location.port: '');
}

document.addEventListener('DOMContentLoaded', setupUI);
function setupUI() {
  listDevices();
  document.querySelector('.nav-toggle').addEventListener ('click', toggleNav);
  document.querySelector('.modal-open').addEventListener('click', openModal);
  document.querySelector('.modal-close').addEventListener('click', closeModal);
  document.querySelector('.modal-background').addEventListener('click', closeModal);
}

function openModal() {
  document.querySelector('.modal').classList.add('is-active');
}

function closeModal() {
  document.querySelector('.modal').classList.remove('is-active');
}

function toggleNav() {
  var nav = document.querySelector(".nav-menu");
  if (nav.classList.contains('is-active')) {
    nav.classList.remove('is-active');
  } else {
    nav.classList.add('is-active');
  }
}

function listDevices() {
  var request = new XMLHttpRequest();
  request.open('GET', location.origin + '/api/devices', true);

  request.onload = function() {
    if (this.status >= 200 && this.status < 400) {
      var devices = JSON.parse(this.response);

      if (devices.length === 0) {
        return;
      }

      var list = document.querySelector('.device-list');
      list.innerHTML = '';

      devices.sort(function(a,b){
        return new Date(b.added) - new Date(a.added);
      });

      devices.forEach(function(l) {

        var t = document.querySelector('.device-template');
        t.content.querySelector('.device-name').textContent = l.name + ' on ' + l.internaladdress;
        //t.content.querySelector('.device-id').textContent = l.id;
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
