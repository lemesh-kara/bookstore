import TopBooksComp from "./TopBooksComp.js";
import AllBooksComp from "./AllBooksComp.js";
import AuthComp from "./AuthComp.js";
import CartComp from "./CartComp.js";
import BookComp from "./BookComp.js";
import ErrorComp from "./ErrorComp.js";
import FAQComp from "./FAQComp.js";
import AboutUsComp from "./AboutUsComp.js";
import ContactsComp from "./ContactsComp.js";
import FeedbackComp from "./FeedbackComp.js";
import NewsComp from "./NewsComp.js";
import constants from "./constants.js";
import Api from "./Api.js";

const routes = [
  { path: "/", component: TopBooksComp },
  { path: "/all-books", component: AllBooksComp },
  { path: "/login", component: AuthComp },
  { path: "/cart", component: CartComp },
  { path: "/faq", component: FAQComp },
  { path: "/contacts", component: ContactsComp },
  { path: "/about-us", component: AboutUsComp },
  { path: "/feedback", component: FeedbackComp },
  { path: "/news", component: NewsComp },
  { path: "/book/:id", component: BookComp, name: "book", props: true },
  { path: "/error/:error", component: ErrorComp, name: "error", props: true },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHashHistory(),
  routes,
});

Vue.createApp({
  data() {
    return {
      query: "",
      limit: 3, // default limit value
      loading: false,
      books: [],
      loaded: false,
    };
  },
  created() {
    this.api = new Api(this.$router);
  },
  watch: {
    $route(to, from) {
      this.query = "";
      this.books = [];
    },
  },
  methods: {
    async searchBooks() {
      if (!this.query) {
        this.books = [];
        return;
      }

      this.loading = true;

      try {
        const response = await axios.get(
          constants.backendRef +
            constants.searchBooksRoute +
            `?query=${this.query}&limit=${this.limit}&distance=3`
        );
        this.books = response.data;
      } catch (error) {
        console.error(error);
      } finally {
        this.loading = false;
      }
    },
    isAuthed() {
      return this.api.isAuthorized();
    },
    signOut() {
      this.api.signOut();
      router.push("/");
      this.$forceUpdate();
    },
  },
})
  .use(router)
  .mount("#app");
