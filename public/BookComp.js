import Api from "./Api.js";
import ToastComp from "./ToastComp.js";

{
  /* <a href="#" @click="openPdf">
<iframe :src="pdfSource" width="50%" height="600px"></iframe>
</a> */
}

export default {
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      book: null,
      reviews: null,
      newReviewText: "",
      newReviewMark: null,
      pdfSource: "",
    };
  },
  created() {
    this.api = new Api(this.$router);
  },
  components: { ToastComp },
  template: `
<ToastComp ref="toast"></ToastComp>
<div class="container">

  <h1>Book</h1>
  <div v-if="book">
    <div class="row">

      <h2 class="mb-4">{{ book.name }}</h2>

      <div class="col-sm-4 mb-3">
        <img :src=book.path_to_cover class="img-fluid m-1" alt="Book Cover">
      </div>

      <div class="col-sm-8 mb-3">
        <p><strong>Author:</strong> {{ book.author }}</p>
        <p><strong>Year:</strong> {{ book.year }}</p>
        <p><strong>Description:</strong> {{ book.description }}</p>
        <p><strong>ISBN:</strong> {{ formatISBN(book.isbn) }}</p>
        <p><strong>Publisher:</strong> {{ book.publisher_data }}</p>
        <p><strong>Review Mark:</strong> {{ Math.round(book.review_mark * 100) / 100 }}</p>
        <p class="text-dark font-weight-bolder h4"><strong>Price:</strong> {{ book.price / 100 }}</p>
        <button class="btn btn-primary mb-3" @click="addToCart(book.id)" :disabled="!isAuthed()">Add to cart</button>
        <br />
        <button class="btn btn-secondary" @click="openPdf">Open preview PDF</button>
      </div>

    </div>

    <br />

    <div v-if="isAuthed()">
      <p class="fs-4">Add review</p>
      <form @submit.prevent="addReview">
        <div class="mb-3">
          <label for="newReviewText" class="form-label">Text:</label>
          <textarea type="text" class="form-control" id="newReviewText" v-model.trim="newReviewText" required></textarea>
        </div>
        <div class="mb-3">
          <label class="form-label d-block">Rating:</label>
          <div class="form-check form-check-inline">
            <input class="form-check-input" type="radio" name="review_mark" id="0" v-model="newReviewMark" :value="0" required>
            <label class="form-check-label" for="0">0</label>
          </div>
          <div class="form-check form-check-inline">
            <input class="form-check-input" type="radio" name="review_mark" id="1" v-model="newReviewMark" :value="1">
            <label class="form-check-label" for="1">1</label>
          </div>
          <div class="form-check form-check-inline">
            <input class="form-check-input" type="radio" name="review_mark" id="2" v-model="newReviewMark" :value="2">
            <label class="form-check-label" for="2">2</label>
          </div>
          <div class="form-check form-check-inline">
            <input class="form-check-input" type="radio" name="review_mark" id="3" v-model="newReviewMark" :value="3">
            <label class="form-check-label" for="3">3</label>
          </div>
          <div class="form-check form-check-inline">
            <input class="form-check-input" type="radio" name="review_mark" id="4" v-model="newReviewMark" :value="4">
            <label class="form-check-label" for="4">4</label>
          </div>
          <div class="form-check form-check-inline">
            <input class="form-check-input" type="radio" name="review_mark" id="5" v-model="newReviewMark" :value="5">
            <label class="form-check-label" for="5">5</label>
          </div>
        </div>
        <div>
          <button type="submit" class="btn btn-primary">Send</button>
        </div>
      </form>
    </div>

    <br />

    <div v-for="review in reviews" :key="review.id">
      <div class="card mb-3">
        <div class="card-body">
          <h5 class="card-title"><strong>Mark:</strong> {{ review.review_mark }}</h5>
          <h6 class="card-subtitle mb-2 text-muted"><strong>User:</strong> {{ review.user.username }}</h6>
          <p class="card-text">{{ review.text }}</p>
        </div>
      </div>
    </div>

  </div>
</div>
      `,
  beforeMount() {
    this.fetchBook();
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
    async fetchBook() {
      this.book = await this.api.fetchBook(this.id);
      this.loadReviews(this.book.id);
      this.pdfSource = this.book.path_to_pdf;
    },
    async addToCart(id) {
      await this.api.addToCart(id);
      this.$refs.toast.showToast();
    },
    async loadReviews(id) {
      this.reviews = await this.api.loadReviews(this.id);
    },
    async addReview() {
      await this.api.addReview(
        this.book.id,
        this.newReviewMark,
        this.newReviewText
      );
      this.newReviewText = "";
      this.newReviewMark = null;
      this.fetchBook();
    },
    openPdf() {
      window.open(this.pdfSource, "_blank", "width=800,height=600");
    },
    formatISBN(number) {
      let formattedNumber = number.replace(
        /(\d{3})(\d{1})(\d{2})(\d{4})(\d{1})/,
        "$1-$2-$3-$4-$5"
      );
      return formattedNumber;
    },
    isAuthed() {
      return this.api.isAuthorized();
    },
  },
};
