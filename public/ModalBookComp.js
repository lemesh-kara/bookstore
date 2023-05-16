import Api from "./Api.js";
import constants from "./constants.js";
import ToastComp from "./ToastComp.js";

{
  /* <a href="#" @click="openPdf">
<iframe :src="pdfSource" width="50%" height="600px"></iframe>
</a> */
}

export default {
  props: {
    book: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      reviews: null,
      newReviewText: "",
      newReviewMark: null,
    };
  },
  created() {
    this.api = new Api(this.$router);
  },
  components: { ToastComp },
  template: `
<button type="button" class="btn border border-primary" data-bs-toggle="modal" data-bs-target="#myModal">
More
</button>

<!-- The Modal -->
<div v-if="book" class="modal fade" id="myModal">
<div class="modal-dialog">
    <ToastComp ref="toast"></ToastComp>
    <div class="modal-content row">

    <!-- Modal Header -->
    <div class="modal-header">
        <h4 class="modal-title">{{ book.name }}</h4>
        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
    </div>

    <!-- Modal body -->
    <div class="modal-body">
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
            <p class="text-dark font-weight-bolder h4"><strong>Price:</strong> $\{{ book.price / 100 }}</p>
            <button class="btn btn-primary mb-1" @click="addToCart(book.id)" :disabled="!isAuthed()">Add to cart</button>
            <br/>
            <button type="button" class="btn btn-secondary mb-1" data-bs-dismiss="modal" @click="goToBook(book.id)">Go to book page</button>
            <br/>
            <button class="btn btn-secondary mb-1" @click="openPdf">Open preview PDF</button>
        </div>
    </div>

    <!-- Modal footer -->
    <div class="modal-footer">
        <button type="button" class="btn btn-danger" data-bs-dismiss="modal">Close</button>
    </div>

    </div>
</div>
</div>
      `,
  mounted() {},
  methods: {
    async addToCart(id) {
      await this.api.addToCart(id);
      this.$refs.toast.showToast();
    },
    openPdf() {
      window.open(this.book.path_to_pdf, "_blank", "width=800,height=600");
    },
    formatISBN(number) {
      let formattedNumber = number.replace(
        /(\d{3})(\d{1})(\d{2})(\d{4})(\d{1})/,
        "$1-$2-$3-$4-$5"
      );
      return formattedNumber;
    },
    goToBook(id) {
      this.$router.push(constants.bookRoute + "/" + id.toString());
    },
    isAuthed() {
      return this.api.isAuthorized();
    },
  },
};
