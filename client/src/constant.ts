
const environment = import.meta.env.VITE_ENVIRONMENT;

let domain: string = "http://localhost:5173/";

if (environment === "production") {
  domain = "https://go-u-sh.vercel.app/";
}

export { domain };
