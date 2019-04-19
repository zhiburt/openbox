import React from "react";
import { Container, Row, Col } from "shards-react";

import Editor from "../components/add-new-post/Editor";



class Viewer extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            body: props.body,
            title: props.title,
          };
    }

    render() {
        return (
        <Container fluid className="main-content-container pt-3 pb-1">
            <Row>
                <Col lg="12" md="12">
                    <Editor body={this.state.body} title={this.state.title}/>
                </Col>
            </Row>
        </Container>
        )
    }
}

export default Viewer;