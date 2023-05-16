import Api from "./Api.js";
import ModalBookComp from "./ModalBookComp.js";
import ToastComp from "./ToastComp.js";

export default {
  data() {
    return {
      topBooks: null,
    };
  },
  created() {
    this.api = new Api(this.$router);
  },
  components: { ModalBookComp, ToastComp },
  template: `
    <ToastComp ref="toast"></ToastComp>
    <div class="container">
    <h1>Top 10 Books</h1>
    <div class="row">
      <div v-for="book in topBooks" :key="book.id" class="col-12 col-md-6 col-lg-5 border border-1 m-1 container">

        <div class="d-flex flex-column h-100">

          <div class="row align-items-center">
            <div class="col-sm-4">
              <img :src=book.path_to_cover class="img-fluid m-1" alt="Book Cover">
            </div>
            <div class="col-sm-8">
              <h2>{{ book.name }}</h2>
              <p class="mb-0"><strong>Author:</strong> {{ book.author }}</p>
              <p class="mb-0"><strong>Year:</strong> {{ book.year }}</p>
              <p class="mb-0"><strong>Description:</strong> {{ book.description }}</p>
              <p class="mb-0"><strong>Review Mark:</strong> {{ Math.round(book.review_mark * 100) / 100 }}</p>
              <p class="mb-0"><router-link :to="{ name: 'book', params: { id: book.id } }">Go to Book</router-link></p>
              <ModalBookComp :book="book"></ModalBookComp>
            </div>
          </div>

          <div class="row justify-content-between mt-auto">
            <span class="col-3 text-dark font-weight-bolder h4"> $\{{ book.price / 100 }}</span>
            <div
              class="col-3 m-1"
              data-bs-toggle="tooltip"
              data-bs-placement="top"
              :title="!isAuthed() ? 'Login to buy books' : 'Click to put to the cart'">
              <button type="button"
                class="btn btn-primary border border-primary w-100"
                :disabled="!isAuthed()"
                @click="addToCart(book.id)"
                >
                <span>Buy</span>
              </button>
            </div>
          </div>

        </div>

      </div>
    </div>
    </div>
      `,
  beforeMount() {
    this.fetchTopBooks();
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
    async fetchTopBooks() {
      this.topBooks = await this.api.fetchTopBooks();
    },
    async addToCart(id) {
      await this.api.addToCart(id);
      this.$refs.toast.showToast();
    },
    isAuthed() {
      return this.api.isAuthorized();
    },
  },
};
