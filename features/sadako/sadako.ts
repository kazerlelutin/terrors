export async function sadako(req: Request) {
  const body = await req.json();
  //TODO 

  /*
  - check AppId and origin ? 
  - use fingerprint on index, for create or update ? 
  - check if the error is already in the database ? 
  - if not, create a new error and send to Discord
  */

  console.log(body);
  return Response.json({
    success: true,
    message: 'Erreur enregistrée',
    timestamp: new Date().toISOString()
  });
}