import constants from "./constants.js";

class Api {
  constructor(router) {
    this.router = router;
    this.client = axios.create({
      baseURL: `${window.location.protocol}//${window.location.hostname}:${window.location.port}`,
    });
    this.client.interceptors.response.use(
      function (response) {
        return response;
      },
      function (error) {
        window.alert(error.response.statusText);
        if(error.response.status === 401) {
          router.push("/login");
        }
        return Promise.reject(error);
      }
    );
  }

  // news
  fetchAllNews() {
    return this.#makeGet(constants.allNewsRoute);
  }

  // Top books
  fetchTopBooks() {
    return this.#makeGet(constants.topBooksRoute);
  }

  // All books
  fetchAllBooks() {
    return this.#makeGet(constants.allBooksRoute);
  }

  // Book
  fetchBook(id) {
    return this.#makeGet(constants.bookRoute + "/" + id);
  }

  // Review
  loadReviews(id) {
    return this.#makeGet(constants.reviewsByBookRoute + id);
  }

  async addReview(bookId, reviewMark, text) {
    if (!(await this.#basicPrecheck())) {
      return null;
    }
    return this.#makeAuthedPost(constants.reviewsRoute, {
      user_id: this.#getCurrentUserId(),
      book_id: bookId,
      review_mark: reviewMark,
      text: text,
    });
  }

  // Cart
  async addToCart(id) {
    if (!(await this.#basicPrecheck())) {
      return null;
    }

    return this.#makeAuthedPost(constants.cartRoute, {
      user_id: this.#getCurrentUserId(),
      book_id: id,
    });
  }

  async getCart() {
    if (!(await this.#basicPrecheck())) {
      return null;
    }

    return this.#makeAuthedGet(
      constants.userCartRoute + this.#getCurrentUserId()
    );
  }

  async removeFromCart(id) {
    if (!(await this.#basicPrecheck())) {
      return null;
    }
    this.#makeAuthedDelete(constants.cartRoute + "/" + id);
  }

  // Auth
  login(username, password) {
    return this.#makePost(constants.loginRoute, {
      username: username,
      password: password,
    });
  }

  signup(email, username, password) {
    return this.#makePost(constants.signUpRoute, {
      email: email,
      username: username,
      password: password,
      role: "user",
    });
  }

  signOut() {
    Cookies.remove(constants.accessTokenName);
    Cookies.remove(constants.refreshTokenName);
  }

  getUsername() {
    if(!this.isAuthorized()) {
      return null;
    }

    return this.#getToken().username;
  }

  isAuthorized() {
    return Cookies.get(constants.accessTokenName) != null;
  }

  // Helpers

  async #basicPrecheck() {
    if (!this.isAuthorized()) {
      this.#goToUnauthorized();
      return false;
    }

    if (this.#isTokenExpired()) {
      await this.#refreshToken();
    }

    return true;
  }

  async #makeGet(path) {
    try {
      const response = await this.client.get(path);
      return response.data;
    } catch (error) {
      console.error(error);
    }
  }

  async #makeAuthedGet(path) {
    try {
      const response = await this.client.get(path, {
        headers: {
          Authorization: "Bearer " + Cookies.get(constants.accessTokenName),
        },
      });
      return response.data;
    } catch (error) {
      console.error(error);
    }
  }

  async #makeAuthedDelete(path) {
    try {
      const access_token = Cookies.get("access_token");
      await this.client.delete(path, {
        headers: { Authorization: "Bearer " + access_token },
      });
    } catch (error) {
      console.error(error);
    }
  }

  async #makePost(path, body) {
    try {
      const response = await this.client.post(path, body);
      return response.data;
    } catch (error) {
      console.error(error);
    }
  }

  async #makeAuthedPost(path, body) {
    try {
      const response = await this.client.post(path, body, {
        headers: {
          Authorization: "Bearer " + Cookies.get(constants.accessTokenName),
        },
      });
      return response.data;
    } catch (error) {
      console.error(error);
    }
  }

  #goToUnauthorized() {
    window.alert("401 Unauthorized\nOops! Try to reload or sign in again!");
  }

  async #refreshToken() {
    try {
      const response = await this.client.get(constants.refreshTokenRoute, {
        headers: {
          Authorization: "Bearer " + Cookies.get(constants.refreshTokenName),
        },
      });
      Cookies.set(constants.accessTokenName, response.data.access_token);
      // Cookies.set(constants.refreshTokenName, response.data.refresh_token);
    } catch (error) {
      console.error(error);
    }
  }

  #isTokenExpired() {
    const accessTokenString = Cookies.get(constants.accessTokenName);
    const accessToken = this.#decodeJWT(accessTokenString);
    const date = new Date(accessToken.exp * 1000);
    const currentDate = new Date();
    return date.getTime() < currentDate.getTime();
  }

  #decodeJWT(token) {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => {
          return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
        })
        .join("")
    );
    return JSON.parse(jsonPayload);
  }

  #getToken() {
    return this.#decodeJWT(Cookies.get(constants.accessTokenName));
  }

  #getCurrentUserId() {
    return this.#getToken().user_id;
  }
}

export default Api;
