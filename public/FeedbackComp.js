import Api from "./Api.js";

export default {
  data() {
    return {
    };
  },
  template: `
<div class = "container">
  <p class="fs-3">Feedback</p>

  <form action="https://formspree.io/f/xvonyljb" method="POST">

    <div class="mb-3 col-5">
      <label for="formControlInput1" class="form-label">Email address (leave empty to stay anonymous):</label>
      <input type="email" class="form-control" id="formControlInput1" placeholder="name@example.com" name="email">
    </div>
    <div class="mb-3">
      <label for="formControlTextarea1" class="form-label">Your feedback:</label>
      <textarea class="form-control" id="formControlTextarea1" rows="3" name="text"></textarea>
    </div>
    <div class="mb-3">
      <button type="submit" class="btn btn-primary">Submit</button>
    </div>
  </form>
</div>
  `,
  created() {
    this.api = new Api(this.$router);
  },
};
