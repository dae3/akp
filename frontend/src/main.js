import './style.css';
import './app.css';
import { Connect, SetSubscription } from '../wailsjs/go/main/App'
import logo from './assets/images/logo-universal.png';

// document.querySelector('#app').innerHTML = ` `;

window.runtime.EventsOn('subscriptions_loaded', (subs) => {
    subs.forEach(sub => document.getElementById('subscriptions').append(new Option(sub.DisplayName, sub.Id)))
    document.getElementById('subscriptions').removeAttribute('disabled')
})

document.getElementById('subscriptions').addEventListener('change', function(e) {
    SetSubscription(this.selectedOptions[0]).then(() => {})
})

window.runtime.EventsOn('resource_change', function(resources) {
    const lis = resources.map(r => {
        const li = document.createElement("li")
        li.append(`${r.Name}: ${r.Status}`)
        return li
    })
    document.getElementById('resources').children[0].replaceChildren(...lis)
})

window.runtime.EventsOn('message', (msg) => {
    document.getElementById('message').replaceChildren(msg)
})

document.getElementById('connect').addEventListener('click', function(e) {
    this.setAttribute('disabled', '')
    Connect().then(() => {})
})
