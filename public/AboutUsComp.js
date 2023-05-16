export default {
    data() {
      return {
        videoUrl: '/video/about-us.mp4',
      };
    },
    template: `
  <div class="container">
    <h2>About Us</h2>
    <div class="row">
      <div class="col-md-6">
        <h3>Our History</h3>
        <p>
          Our bookstore was founded in 1990 by John Smith, a lifelong book lover and avid reader. John had a passion for books and wanted to share that passion with others, so he decided to open a bookstore in his hometown of Springfield, USA.
        </p>
        <p>
          In the early years, the store was small but quickly gained a loyal following among book lovers in the area. As the store grew, John added more inventory and expanded the store's offerings to include not just books, but also magazines, journals, and other reading materials.
        </p>
        <p>
          Over the years, our bookstore has become a beloved institution in the community, known for our wide selection of books, knowledgeable staff, and commitment to providing an exceptional customer experience.
        </p>
      </div>
      <div class="col-md-6">
        <video controls class="img-fluid">
          <source v-bind:src="videoUrl" type="video/mp4">
          Your browser does not support the video tag.
        </video>
      </div>
    </div>
  </div>
          `,
  };
  