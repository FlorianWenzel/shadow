import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'

import App from './App.vue'
import router from './router'

import './assets/theme.scss'
import '@fortawesome/fontawesome-free/css/all.css'
import 'primeflex/primeflex.css'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import Toolbar from 'primevue/toolbar'

import 'chart.js'
import Card from 'primevue/card'
import 'chartjs-chart-geo'

const app = createApp(App)

app.component('Button', Button)
app.component('InputText', InputText)
app.component('IconField', IconField)
app.component('InputIcon', InputIcon)
app.component('Toolbar', Toolbar)
app.component('Card', Card)

app.use(PrimeVue)
app.use(createPinia())
app.use(router)

app.mount('#app')
