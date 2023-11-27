<template>
  <div class="container">

    <section class="section">
      <div class="block">
        <form method="get">
          <b-field>
            <b-input name="s" placeholder="Search..." v-model="searchText" type="search"></b-input>
            <p class="control">
              <b-button native-type="submit" class="button is-primary">Search</b-button>
              <b-button v-if="searchText.length > 0" class="button" @click="handleClear">Clear</b-button>
            </p>
          </b-field>
        </form>
      </div>

      <div class="block">
        <b-table :data="mails"
                 :striped="true"
                 :hoverable="true"
                 :per-page="pageSize"
                 :pagination-simple="true"
                 :loading="loading"
                 :total="total"
                 paginated
                 backend-pagination
                 aria-next-label="Next page"
                 aria-previous-label="Previous page"
                 @page-change="onPageChange"
                 :current-page="currentPage">
          <template #cell(subject)="row">
            <div class="block">
              {{ row.value }}
            </div>
          </template>

          <b-table-column label="ID" v-slot="props">{{props.row.id}}</b-table-column>
          <b-table-column label="Subject" v-slot="props" width="300">
                <div class="block is-align-content-center">
                  <a @click="callPreview(props.row.id)" class="subtitle link"> {{ props.row.subject }} </a>
                </div>
          </b-table-column>
          <b-table-column label="From" v-slot="props">{{props.row.from}}</b-table-column>
          <b-table-column label="To" v-slot="props">{{props.row.to}}</b-table-column>
          <b-table-column label="Date" v-slot="props">{{props.row.date}}</b-table-column>
          <b-table-column label="Create Time" v-slot="props">{{props.row.created_at}}</b-table-column>

          <template #empty>
            <div class="has-text-centered">No records</div>
          </template>
        </b-table>
      </div>

      <b-modal
          v-model="showPreview"
          has-modal-card
          aria-role="dialog"
          aria-label="Email Preview"
          close-button-aria-label="Close"
          aria-modal
          :width="750"
          scroll="keep">
        <preview :mail="previewData"></preview>
      </b-modal>

    </section>
  </div>
</template>

<script>
import axios from "axios";
import Preview from "./preview.vue";

const http = axios.create({
  baseURL: global.API_URI || "/api/"
})

http.interceptors.response.use(function (resp) {
  let json = resp.data;
  if(json && json.status) {
    return Promise.resolve(json.data);
  }
  if(json && json.error) {
    return Promise.reject(new Error(json.error));
  }
  return Promise.reject(new Error('api error'));
}, function (error) {
  let resp = error.response.data;
  if (resp.error) {
    return Promise.reject(new Error(resp.error));
  }
  return Promise.reject(new Error('api error'));
});

export default {
  name: "App",

  components: {
    Preview,
  },

  data() {
    return {
      previewData:{},
      showPreview: false,
      mails: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      loading: false,
    }
  },

  computed:{
    columns() {
      return [
        {
          field: 'id',
          label: 'ID',
        },
        {
          field: 'subject',
          label: 'Subject',
        },
        {
          field: 'from',
          label: 'From',
        },
        {
          field: 'to',
          label: 'To',
        },
        {
          field: 'date',
          label: 'Date',
        },
        {
          field: 'created_at',
          label: 'Created Time',
        },
      ];
    },
    searchText() {
      const query = new URLSearchParams(location.search);
      return query.get("s") || "";
    },
  },

  mounted() {
    this.fetchData();
  },

  methods: {
    fetchData() {
      this.loading = true;
      http.get('/mail', {
        params: {
          pageSize: this.pageSize,
          page: this.currentPage,
          s: this.searchText,
        }
      })
      .then(data => {
        this.mails = data.items || [];
        this.pageSize = data.pagination?.pageSize || this.pageSize;
        this.currentPage = data.pagination?.current || this.currentPage;
        this.total = data.pagination?.total || this.total;
      })
      .finally(() => {
        this.loading = false;
      })
    },
    handleClear() {
      location.href = location.pathname;
    },
    onPageChange(page) {
      this.currentPage = page;
      this.fetchData();
    },
    callPreview(id) {
      this.loading = true;
      http.get(`/mail/${id}`)
          .then(data => {
            this.previewData = data;
            this.showPreview = true;
          })
          .finally(() => {
            this.loading = false;
          })
    }
  }
}
</script>
