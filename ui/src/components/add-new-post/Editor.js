import React from "react";
import ReactQuill from "react-quill";
import { Card, CardBody, Form } from "shards-react";

import "react-quill/dist/quill.snow.css";
import "../../assets/quill.css";



const Editor = (props) => (
  <Card small className="mb-3">
    <CardBody>
      <Form className="add-new-post">
        <ReactQuill value={props.body} className="add-new-post__editor mb-1" />
        {props.title} 
      </Form>
    </CardBody>
  </Card>
);

export default Editor;