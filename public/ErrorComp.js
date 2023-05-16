import constants from "./constants.js";

export default {
  props: {
    error: {
      type: String,
      required: true,
    },
  },
  template: `
      <div>
        <h1>{{ error }}</h1>
        <h2>Something bad happened... Try to reload or sign in once again!</h2>
      </div>
      `,
};
