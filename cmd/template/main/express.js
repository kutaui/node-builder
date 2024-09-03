import express from 'express';

const port = process.env.PORT || 5000;

const app = express();


app.get('/', (req, res) => {
    res.send('Hello World');
})


app.listen(port, () => {
    console.log(`Server started at ${port}`);
})