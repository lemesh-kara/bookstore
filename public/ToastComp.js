export default {
  props: {
    message: {
      type: String,
      default: "Book added to the shopping cart!",
    },
  },
  template: `
<div class="toast align-items-center position-fixed top-0 end-0" role="alert" aria-live="assertive" aria-atomic="true">
  <div class="d-flex">
    <div class="toast-body">
    {{ message }}
    </div>
    <button type="button" class="btn-close me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
  </div>
</div>
      `,
  methods: {
    showToast() {
      new bootstrap.Toast(document.querySelector(".toast")).show();
    },
  },
};
