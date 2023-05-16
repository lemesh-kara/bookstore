import Api from "./Api.js";

export default {
  data() {
    return {
      username: "",
      password: "",
      email: "",
      isLogin: true,
    };
  },
  created() {
    this.api = new Api(this.$router);
  },
  template: `
    <div class="container">
      <ul class="nav nav-tabs">
        <li class="nav-item">
          <a class="nav-link" :class="{ active: isLogin }" @click="() => { this.isLogin = true; }">Login</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" :class="{ active: !isLogin }" @click="() => { this.isLogin = false; }">Sign Up</a>
        </li>
      </ul>
      <br/>
      <div class="tab-content">
        <div class="tab-pane fade" :class="{ show: isLogin, active: isLogin }">
          <form @submit.prevent="login">
            <div class="mb-3 col-sm-6">
              <label for="username" class="form-label">Username:</label>
              <input class="form-control" type="text" id="username" v-model.trim="username" required />
            </div>
            <div class="mb-3 col-sm-6">
              <label for="password" class="form-label">Password:</label>
              <input class="form-control" type="password" id="password" v-model="password" required />
            </div>
            <div class="mb-3">
              <button type="submit" class="btn btn-primary">Login</button>
            </div>
          </form>
        </div>

        <div class="tab-pane fade" :class="{ show: !isLogin, active: !isLogin }">
          <form @submit.prevent="signup">
            <div class="mb-3 col-sm-6">
              <label for="email" class="form-label">Email:</label>
              <input class="form-control" type="email" id="email" v-model.trim="email" required />
            </div>
            <div class="mb-3 col-sm-6">
              <label for="usernameSignUp" class="form-label">Username:</label>
              <input class="form-control" type="text" id="usernameSignUp" v-model.trim="username" required />
            </div>
            <div class="mb-3 col-sm-6">
              <label for="passwordSignUp" class="form-label">Password:</label>
              <input class="form-control" type="password" id="passwordSignUp" v-model="password" required />
            </div>
            <div class="mb-3">
              <button type="submit" class="btn btn-primary">Sign Up</button>
            </div>
          </form>
        </div>
      </div>
    </div>
      `,
  mounted() {},
  methods: {
    async login() {
      const response = await this.api.login(this.username, this.password);
      Cookies.set("access_token", response.access_token);
      Cookies.set("refresh_token", response.refresh_token);
      this.username = "";
      this.password = "";
      this.$router.push("/");
    },

    async signup() {
      const response = await this.api.signup(
        this.email,
        this.username,
        this.password
      );
      Cookies.set("access_token", response.access_token);
      Cookies.set("refresh_token", response.refresh_token);
      this.username = "";
      this.email = "";
      this.password = "";
      this.$router.push("/");
    },
  },
};
