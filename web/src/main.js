import Vue  from 'vue';
import App from './App.vue';
import { Table, Input, Field, Button } from 'buefy'
import 'buefy/dist/buefy.css';
import '@mdi/font/css/materialdesignicons.min.css'

Vue.use(Table);
Vue.use(Input);
Vue.use(Button);
Vue.use(Field);

new Vue({
  render: h => h(App),
}).$mount('#app');
