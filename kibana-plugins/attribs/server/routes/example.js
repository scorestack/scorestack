export default function(server) {
  server.route({
    path: '/api/attribs/example',
    method: 'GET',
    handler() {
      return { time: new Date().toISOString() };
    },
  });
}
