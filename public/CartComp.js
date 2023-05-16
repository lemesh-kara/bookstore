import helpers from "./helpers.js";
import Api from "./Api.js";

export default {
  data() {
    return {
      userCart: null,
    };
  },
  created() {
    this.api = new Api(this.$router);
  },
  template: `
      <div>
        <h1 class="mb-4">Cart</h1>
        <div v-if="!userCart || userCart.length === 0">
          <p>Nothing in cart yet</p>
        </div>
        <div v-else v-for="cartRecord in userCart" :key="cartRecord.id" class="card mb-3">
          <div class="card-body">
            <h5 class="card-title">{{ cartRecord.book.name }}</h5>
            <p class="card-text">Price: {{ cartRecord.book.price.toFixed(2) / 100 }}</p>
            <p class="card-text"><router-link :to="{ name: 'book', params: { id: cartRecord.book.id } }">Go to Book</router-link></p>
            <button @click="removeCart(cartRecord.id)" class="btn btn-danger">Remove</button>
          </div>
        </div>
        <div>
          <h2>Total Amount: {{ totalAmount.toFixed(2) / 100 }}</h2>
        </div>
        <br/>
        <button disabled=true class="btn btn-primary" data-toggle="popover" title="In development warning" data-content="Still in development... Sorry">Make order</button>
      </div>
      `,
  mounted() {
    this.fetchUserCart();
  },
  computed: {
    totalAmount: function () {
      if (!this.userCart || this.userCart.length === 0) {
        return 0;
      }
      return this.userCart.reduce(function (acc, cartRecord) {
        return acc + cartRecord.book.price;
      }, 0);
    },
  },
  methods: {
    async fetchUserCart() {
      this.userCart = await this.api.getCart();
    },
    async removeCart(id) {
      await this.api.removeFromCart(id);
      this.userCart.splice(
        this.userCart.findIndex((cartRecord) => cartRecord.id === id),
        1
      );
    },
  },
};
