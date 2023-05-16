export default {
  data() {
    return {
      activeIndex: null,
      faqs: [
        {
          id: 0,
          question: "How can I place an order?",
          answer:
            "To place an order, simply browse our selection of books and add the ones you want to your shopping cart. When you're ready to check out, click on the cart icon at the top of the page and follow the instructions to complete your purchase.",
        },
        {
          id: 1,
          question: "What payment methods do you accept?",
          answer:
            "We accept all major credit cards, including Visa, MasterCard, American Express, and Discover. You can also use PayPal to complete your purchase.",
        },
        {
          id: 2,
          question: "How long does shipping take?",
          answer:
            "Shipping times vary depending on your location and the shipping method you choose at checkout. Standard shipping usually takes 3-5 business days, while expedited shipping takes 1-2 business days.",
        },
        {
          id: 3,
          question: "What is your return policy?",
          answer:
            "If you're not satisfied with your purchase, you can return it within 30 days for a full refund. Just contact us to initiate the return process.",
        },
      ],
    };
  },
  template: `
  <div class="accordion accordion-flush" id="accordionFlush">
  <div class="accordion-item"  v-for="faq in faqs" :key="faq.id">
    <h2 class="accordion-header" :id="'flush-heading' + faq.id">
      <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" :data-bs-target="'#flush-collapse' + faq.id" aria-expanded="false" :aria-controls="'flush-collapse' + faq.id">
        {{ faq.question }}
      </button>
    </h2>
    <div :id="'flush-collapse' + faq.id" class="accordion-collapse collapse" :aria-labelledby="'flush-heading' + faq.id" data-bs-parent="#accordionFlush">
      <div class="accordion-body">{{ faq.answer }}</div>
    </div>
  </div>
</div>
        `,
};
