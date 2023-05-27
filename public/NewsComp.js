import Api from "./Api.js";
import ModelNewsComp from "./ModelNewsComp.js";
import ToastComp from "./ToastComp.js";

export default {
  data() {
    return {
      allNews: null,
    };
  },
  created() {
    this.api = new Api(this.$router);
  },
  components: { ModelNewsComp, ToastComp },
  template: `
<ToastComp ref="toast"></ToastComp>
<div class="container">
<h1>News</h1>
<div class="row">
  <div v-for="news in allNews" :key="news.id" class="col-12 border border-1 m-1 container">

    <div class="d-flex flex-column h-100">

      <div class="row align-items-center">
        <div class="col-sm-4">
          <img :src=news.path_to_picture class="img-fluid m-1" alt="News Cover">
        </div>
        <div class="col-sm-8">
          <h2>{{ news.title }}</h2>
          <p class="mb-0"><strong>Text:</strong> {{ truncateString(news.text) }}</p>
          <ModelNewsComp :news="news"></ModelNewsComp>
        </div>
      </div>

    </div>

  </div>
</div>
</div>
      `,
  beforeMount() {
    this.fetchAll();
  },
  mounted() {
    var tooltipTriggerList = [].slice.call(
      document.querySelectorAll('[data-bs-toggle="tooltip"]')
    );
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
      return new bootstrap.Tooltip(tooltipTriggerEl);
    });
  },
  methods: {
    async fetchAll() {
      this.allNews = await this.api.fetchAllNews();
    },
    async addToCart(id) {
      await this.api.addToCart(id);
      this.$refs.toast.showToast();
    },
    isAuthed() {
      return this.api.isAuthorized();
    },
    truncateString(text) {
      if (text.length <= 80) {
        return text;
      } else {
        return text.slice(0, 80) + "...";
      }
    },
  },
};
