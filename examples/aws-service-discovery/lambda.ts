import { Resource } from "neo";

export async function handler() {
  const reponse = await fetch(
    `http://${Resource.MyService.service}`
  );

  return {
    statusCode: 200,
    body: await reponse.text(),
  };
}
