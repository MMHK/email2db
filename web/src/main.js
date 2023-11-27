import Vue  from 'vue';
import App from './App.vue';
import { Table, Input, Field, Button, Modal } from 'buefy'
import 'buefy/dist/buefy.css';
import '@mdi/font/css/materialdesignicons.min.css'

Vue.use(Table);
Vue.use(Input);
Vue.use(Button);
Vue.use(Field);
Vue.use(Modal);

new Vue({
  render: h => h(App),
}).$mount('#app');
