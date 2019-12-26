/* eslint jsx-a11y/anchor-is-valid: 0 */

import React from "react";
import {
  Container,
  Row,
  Col,
  Card,
  CardBody,
  CardHeader,
  CardFooter,
  Button,
  Modal,
  Breadcrumb,
  BreadcrumbItem,
  InputGroup,
  FormInput,
  InputGroupAddon,
} from "shards-react";

import PageTitle from "../components/common/PageTitle";
import Toolbar from "../components/common/Toolbar";
import Viewer from "./Viewer"
import ButtonFileUpload from "../components/components-overview/ButtonFileUpload";

class BlogPosts extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      error: null,
      isLoaded: false,
      files: [],
      recievedFiles: [],
      breadcrumds: ["~"],
      indexFile: -1,
      openModal: false,
    };

    this.toggle = this.toggle.bind(this);
    this.pushFile = this.pushFile.bind(this);
    this.buildJsonFile = this.buildJsonFile.bind(this);
    this.createFolder = this.createFolder.bind(this);
  }

  toggle() {
    this.setState({
      openModal: !this.state.openModal
    });
  }

  handleClick(index, e) {
    console.log(index, this.state.files[index]);
    if (this.state.files[index].is_folder == true) {
      var name = this.state.files[index].name
      this.state.breadcrumds.push(name)

      this.setState({ files: this.state.files[index].files })
      return
    }

    this.setState({ indexFile: index });
    this.toggle();
  }

  _renderSubComp() {
    if (this.state.indexFile == -1) {
      return
    }

    console.log("BODY", this.state.files[this.state.indexFile])

    const { openModal } = this.state;
    return (
      <Modal open={openModal} size="lg" onClick={this.toggle}>
        <Viewer body={this.state.files[this.state.indexFile].body} title={this.state.files[this.state.indexFile].name} />
        <Button squared theme="info" onClick={this.toggle}>Close</Button>
      </Modal>)
  }

  componentDidMount() {
    fetch('http://localhost:8082/files/owner/zhiburt', {
      method: "GET",
    }).then(res => res.json())
      .then(
        (result) => {
          console.log("+++++++++", result["files"]);
          this.setState({
            isLoaded: true,
            files: result["files"],
            recievedFiles: result["files"],
          });
        },
        (error) => {
          console.log("====", error)
        }
      )
  }

  e(e) {
    console.log("rerender");
  }

  breadClickHandler(nameFolder, index, deep = 1, tempFiles = null, e) {
    if (tempFiles == null) {
      tempFiles = this.state.recievedFiles
      deep = 1
    }

    console.log("index ---")
    console.log(index)
    console.log(nameFolder)
    console.log("index ---")

    if (index == 0) {
      this.setState({ files: tempFiles, breadcrumds: ["~"], indexFile: -1, });
      return
    }

    console.log("deep " + deep)
    var f = getFile(tempFiles, this.state.breadcrumds[deep]);
    if (f == null) {
      console.log("ERORR F == null")
      return
    } else if (index == deep && nameFolder == f.name) {
      this.setState({ files: f.files, breadcrumds: this.state.breadcrumds.slice(0, deep + 1), indexFile: -1, });
      return;
    } else {
      this.breadClickHandler(nameFolder, index, deep + 1, f.files, e)
    }

    //if nameFold == breadName && deep == index
    //return
    //breadName 1
    //go to file with name of bread
    //getBreadName2
    //go ...
  }

  buildJsonFile(file) {
    if (this.state.recievedFiles == this.state.files || this.state.recievedFiles.length == 0) {
      return file;
    }

    var files = this.state.recievedFiles;
    var json = null;
    for (var i = 1; i < this.state.breadcrumds.length; i++) {
      var f = getFile(files, this.state.breadcrumds[i]);
      files = f.files;
      f.files = null;

      console.log("JSON-->>", json)

      if (json == null) {
        json = f
      } else {
        json.files = [f]
      }

      if (i + 1 == this.state.breadcrumds.length) {
        json.files = [file]
      }
    }

    return json
  }

  pushFile(file, url) {
    const owner_id = "zhiburt"
    var d = this.buildJsonFile({ "name": file.name, "body": (url), "owner_id": owner_id })

    console.log("hello world JSON", d);
    console.log("hello world BODY" + url);

    fetch('http://localhost:8082/files', {
      mode: "no-cors",
      method: "POST",
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(d)
    })
    // invalid responce currently
    .then((result) => result.json())
    .then((result) => {
      console.log("************");
      console.log(result);
    })
    .catch((error) => {
      this.componentDidMount();
      console.error(error);
    });
  }

  createFolder(file) {
    const owner_id = "zhiburt"

    if (this.state.dictionary_name == "") {
      this.state.dictionary_name = "please write data";
      return
    }

    var d = this.buildJsonFile({ "name": this.state.dictionary_name, "is_folder": true, "owner_id": owner_id, })

    console.log("createFolder hello world JSON", d);

    fetch('http://localhost:8082/files', {
      mode: "no-cors",
      method: "POST",
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(d)
    }).then(res => res.json())
      .then(
        (result) => {
          console.log("RESULT S2222", result);
        },
        (error) => {
          console.log("ERORR", error)
        }
      )
  }

  render() {
    const { files, breadcrumds } = this.state;
    return (
      <Container fluid className="main-content-container px-4">
        {/* Page Header */}
        <Row noGutters className="page-header py-4">
          <PageTitle sm="4" title="Blog Posts" subtitle="Components" className="text-sm-left" />
        </Row>

        <Row noGutters>
          <Col lg="8">
            <Breadcrumb>
              {breadcrumds.map((bread, idx) => (
                <BreadcrumbItem active key={"bread_" + idx}>
                  <a href="#" onClick={this.breadClickHandler.bind(this, bread, idx)}>{bread}</a>
                </BreadcrumbItem>
              ))}
            </Breadcrumb>


            <Row>

              {/* Third Row of Posts */}
              {this.state.files == undefined ? null : this.state.files.map((post, idx) => (
                <Col lg="3" key={idx}>
                  <Card small className="card-post mb-4" onClick={this.handleClick.bind(this, idx)}>
                    {
                      post.is_folder != true ? "" :
                      <CardHeader small>
                        <i className="far fa-folder mr-1" />
                      </CardHeader>
                    }
                    <CardBody>
                      <h5 className="card-title">{post.name}</h5>
                      <p className="card-text text-muted">Some description</p>
                    </CardBody>
                    <CardFooter className="border-top d-flex">
                      <div className="card-post__author d-flex">
                        <a
                          href="#"
                          className="card-post__author-avatar card-post__author-avatar--small"
                          style={{ backgroundImage: `url('${require("../images/avatars/3.jpg")}')` }}
                        >
                          Written by James Khan
                    </a>
                        <div className="d-flex flex-column justify-content-center ml-3">
                          <span className="card-post__author-name">
                            {post.author}
                          </span>
                          <small className="text-muted">04.14.2019</small>
                        </div>
                      </div>
                      <div className="my-auto ml-auto">
                        <Button size="sm" theme="" circle>
                          <i className="far fa-bookmark mr-1" />
                        </Button>
                      </div>
                    </CardFooter>
                  </Card>
                </Col>
              ))}
            </Row>
          </Col>
          <Col lg="4" className="pl-3">
            <Card small className="mb-3">
              <CardHeader className="border-bottom">
                <h6 className="m-0">Settings</h6>
              </CardHeader>
              <CardBody className="p-0">
                <div className="px-3 pt-2">
                  <ButtonFileUpload callBack={this.pushFile} />
                </div>
                <InputGroup className="pb-2 px-3">
                  <FormInput placeholder="New folder" value={this.state.dictionary_name} onInput={e => this.setState({ "dictionary_name": e.target.value })} />
                  <InputGroupAddon type="append">
                    <Button theme="white" className="px-2" onClick={e => this.createFolder()} >
                      <i className="material-icons">add</i>
                    </Button>
                  </InputGroupAddon>
                </InputGroup>
                <div className="pb-3 px-3" >
                  <Toolbar callBack={this.e} />
                </div>
              </CardBody>
            </Card>
          </Col>
        </Row>

        {this._renderSubComp()}

      </Container>
    );
  }
}

function getFile(files, name) {
  for (var i = 0; i < files.length; i++) {
    if (files[i].name == name) {
      console.log("name " + i + " " + files[i].name);
      return files[i];
    }
  }

  return null;
}

export default BlogPosts;
