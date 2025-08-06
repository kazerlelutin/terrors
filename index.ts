import { sadako } from "./features/sadako/sadako";

Bun.serve({
  port: Bun.env.PORT || 3000,
  // `routes` requires Bun v1.2.3+
  routes: {
    "/": async () => {

      return Response.json({
        message: "Hello World"
      });
    },
    // Endpoint API pour recevoir les erreurs
    "/sadako": {
      POST: sadako,
      OPTIONS: () => {
        return new Response(null, {
          headers: {
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "POST, OPTIONS",
            "Access-Control-Allow-Headers": "Content-Type"
          }
        });
      }
    },

    // Static routes
    "/cdn/terrors.js": new Response(await Bun.file("features/terrors/terrors.js").text(), {
      headers: {
        "Content-Type": "application/javascript; charset=utf-8",
        "Access-Control-Allow-Origin": "*"
      }
    }),

    "/*": new Response("Not Found", { status: 404 }),
  },
});

console.log(`Server is running on http://localhost:${Bun.env.PORT || 3000}`);