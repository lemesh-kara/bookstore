const backendRef = `${window.location.protocol}//${window.location.hostname}:${window.location.port}`;
const bookRoute = "/book";
const topBooksRoute = "/book?limit=10&field=review_mark&asc=false";
const searchBooksRoute = "/book/search";
const allBooksRoute = "/book/all";
const loginRoute = "/login";
const signUpRoute = "/signup";
const cartRoute = "/cart";
const userCartRoute = "/cart/searchby/user?userid=";
const reviewsRoute = "/review";
const reviewsByBookRoute = "/review/searchby/book?bookid=";
const reviewsByUserRoute = "/review/searchby/book?userid=";
const refreshTokenRoute = "/user/refresh";

const errorRoute = "/error";

const accessTokenName = "access_token";
const refreshTokenName = "refresh_token";

const allNewsRoute = "/news/all";
const newsRoute = "/news/";

export default {
  backendRef,
  bookRoute,
  topBooksRoute,
  searchBooksRoute,
  allBooksRoute,
  loginRoute,
  signUpRoute,
  cartRoute,
  userCartRoute,
  reviewsRoute,
  reviewsByBookRoute,
  reviewsByUserRoute,
  refreshTokenRoute,

  errorRoute,

  accessTokenName,
  refreshTokenName,

  allNewsRoute,
  newsRoute,
};
