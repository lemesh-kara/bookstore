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
    news: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
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
<div v-if="news" class="modal fade" id="myModal">
<div class="modal-dialog">
    <ToastComp ref="toast"></ToastComp>
    <div class="modal-content row">

    <!-- Modal Header -->
    <div class="modal-header">
        <h4 class="modal-title">{{ news.Title }}</h4>
        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
    </div>

    <!-- Modal body -->
    <div class="modal-body">
        <div class="col-sm-4 mb-3">
            <img :src=news.path_to_picture class="img-fluid m-1" alt="News Cover">
        </div>

        <div class="col-sm-8 mb-3">
            <p class="mb-0"><strong>Text:</strong> {{ news.text }}</p>
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
    goToBook(id) {
      this.$router.push(constants.bookRoute + "/" + id.toString());
    },
    isAuthed() {
      return this.api.isAuthorized();
    },
  },
};
