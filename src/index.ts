import querystring from 'querystring'
import express from "express";
import cors from "cors"
import pino from "pino-http"
import bodyparser from 'body-parser';

import { ENV } from './env';
import axios from "axios";
import { z } from 'zod'

const app = express();

const IdentitySchema = z.object({
  firstName: z.string(),
  lastName: z.string(),
  middleName: z.string().optional(),
  dateOfBirth: z.iso.date(),
})

app.use(pino())
app.use(cors())
app.use(bodyparser.json())

app.get('/health', (req, res) => {
  req.log.info('health');
  res.status(200).send('OK')
})

app.get('/enrollment', async (req, res) => {
  try {
    const identity = IdentitySchema.parse(req.body);

    req.log.info({body: identity}, "got request");

    const token = await axios.post(ENV.NSC_ACCESS_TOKEN_URL, querystring.stringify({
      client_id: ENV.NSC_CLIENT_ID,
      client_secret: ENV.NSC_CLIENT_SECRET,
      grant_type: "client_credentials",
      scope: "vs.api.insights",
    }))

    const { access_token } = token.data;

    const enrollment = await axios.post(`${ENV.NSC_BASE_URL}/insights/v3/a2/submit-request`,
      JSON.stringify({
        firstName: identity.firstName,
        lastName: identity.lastName,
        dateOfBirth: identity.dateOfBirth,
        "accountId": ENV.NSC_ACCOUNT_ID,
        "terms": "y",
        "endClient": "CMS"
      }),
      {
        headers: {
          'Content-Type': "application/json",
          'Authorization': `Bearer ${access_token}`
        }
      }
    )

    res.send(enrollment.data)
  } catch (e) {
    req.log.error(e)
    res.send({ error: String(e) })
  }
})

app.get('/', (_req, res) => {
  res.send('Hello, World!')
})

app.listen(ENV.PORT, (e) => {
  if (e) {
    console.error(e)
  } else {
    console.log(`listening on port ${ENV.PORT}`);
  }
})
