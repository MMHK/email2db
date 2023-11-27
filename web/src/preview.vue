<template>
  <div class="modal-card">
    <header class="modal-card-head">
      <p class="subtitle">{{mail.subject}}</p>
    </header>
    <section>
      <p class="modal-card-body"><strong>From: </strong> {{mail.from}} <br> <strong>To: </strong>{{mail.to.join(",")}}</p>
    </section>
    <section class="modal-card-body">
      <div class="card">
        <div class="card-content">
          <div class="content" v-html="previewHTML">
          </div>
        </div>
      </div>
    </section>
    <section class="modal-card-body" v-if="mail.attachments && mail.attachments.length > 0">
      <h2 class="subtitle">Attachments:</h2>
      <ul class="columns">
        <li class="column" v-for="attachment in mail.attachments"><a
            :href="attachmentURL(attachment.id)" target="_blank"
            class="button is-primary">{{attachment.name}}</a></li>
      </ul>
    </section>
    <footer class="modal-card-foot is-justify-content-flex-end">
      <b-button label="Close" @click="$parent.close()" />
    </footer>
  </div>
</template>
<script>
import axios from "axios";
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

export default  {
  name: "preview",

  props:{
    mail: {
      subject: "",
      from: "",
      to: [],
      html: "",
      attachments: [],
    }
  },

  computed:{
    previewHTML() {
      let html = this.mail?.html || "";

      if (this.mail?.attachments && this.mail?.attachments.length > 0) {
        html = html.replace(/src="cid:([^"]+)"/ig, (match, p1) => {
          return `src="${this.viewAttachmentImage(p1)}"`;
        })
      }

      return html;
    }
  },

  methods: {
    attachmentURL(id) {
      const baseURL = global.API_URI || "/api/";
      return `${baseURL}attachment/${id}`
    },
    viewAttachmentImage(cid) {
      const baseURL = global.API_URI || "/api/";
      const attachment_id = this.mail?.attachments.find((row) => {
        return row["content-id"] == cid;
      })?.id;
      return attachment_id ? `${baseURL}attachment/${attachment_id}` : cid;
    }
  }
}
</script>
