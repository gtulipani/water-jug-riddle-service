import React, { Component } from "react";
import axios from "axios";
import { Button, Form, Header, Icon, Input, Message, Table } from "semantic-ui-react";

let endpoint = "http://localhost:8080";

class Riddle extends Component {
    constructor(props) {
        super(props);
    
        this.state = {
          params: {
            x: "",
            y: "",
            z: ""
          },
          resolution: {
            filled: false,
            loading: false,
            ok: false,
            jug: "",
            operations: [],
            total_steps: "",
            error: {
              message: "",
              description: ""
            }
          }
        };
    }

    onSubmit = () => {
        let { x, y, z } = this.state;
        let scope = this
        if ((x !== undefined) && (y !== undefined) && (z !== undefined)){
          scope.setState({
            resolution: {
              filled: false,
              loading: true,
              ok: false,
            }
          })
          axios
            .get(
              endpoint + "/api/v1/riddle?x=" + x + "&y=" + y + "&z=" + z
            )
            .then((res) => {
              scope.setState({
                resolution: {
                  jug: res.data.jug,
                  operations: res.data.operations,
                  total_steps: res.data.total_steps,
                  filled: true,
                  loading: false,
                  ok: true,
                }
              })
            })
            .catch(function (error) {
              console.log(error);
              if(error.response) {
                console.log(error.response);
              }
              scope.setState({
                resolution: {
                  filled: true,
                  loading: false,
                  ok: false,
                  error: {
                    message: error.response.data.message,
                    description: error.response.data.description,
                  }
                }
              })
            });
        }
      };

    onChangeX = (event) => {
        this.setState({
          x: event.target.value,
          y: this.state.y,
          z: this.state.z
        });
      };
    
    onChangeY = (event) => {
        this.setState({
            x: this.state.x,
            y: event.target.value,
            z: this.state.z
        });
    };

    onChangeZ = (event) => {
        this.setState({
            x: this.state.x,
            y: this.state.y,
            z: event.target.value
        });
    };
    
    render() {
        return (
        <div>
            <div className="row">
            <Header className="header" as="h2">
                Water Jug Riddle
            </Header>
            </div>
            <div className="row">
            <Form>
                <Input
                type="text"
                name="x"
                onChange={this.onChangeX}
                value={this.state.x}
                fluid
                placeholder="X"
                />
                <Input
                type="text"
                name="y"
                onChange={this.onChangeY}
                value={this.state.y}
                fluid
                placeholder="Y"
                />
                <Input
                type="text"
                name="y"
                onChange={this.onChangeZ}
                value={this.state.z}
                fluid
                placeholder="Z"
                />
              <Button
                name="Calculate"
                onClick={this.onSubmit}
              >
                Calculate Operations
              </Button>
            </Form>
            </div>
            {this.state.resolution.filled && this.state.resolution.ok &&
              <div>
              <Table>
                <Table.Header>
                  <Table.Row>
                  <Table.HeaderCell>Step Number</Table.HeaderCell>
                    <Table.HeaderCell>Operation Type</Table.HeaderCell>
                    <Table.HeaderCell>Jugs</Table.HeaderCell>
                    <Table.HeaderCell>Amount of Water</Table.HeaderCell>
                  </Table.Row>
                </Table.Header>

                <Table.Body>
                {this.state.resolution.operations.map(op => {
                  return (
                    <Table.Row>
                      <Table.Cell>{"#" + op.step}</Table.Cell>
                      <Table.Cell>{op.operation.toUpperCase()}</Table.Cell>
                      <Table.Cell>{op.jug !== undefined ? op.jug.toUpperCase() : op.jug_origin.toUpperCase() + " to " + op.jug_destination.toUpperCase()}</Table.Cell>
                      <Table.Cell>{op.amount}</Table.Cell>
                    </Table.Row>
                  )
                })}
              </Table.Body>
              </Table>
              <Header className="header" as="h4">
                Jug with desired amount of water: {this.state.resolution.jug.toUpperCase()}
              </Header>
              </div>
            }
            {this.state.resolution.loading &&
              <div>
                <Header className="header" as="h4">
                Calculating!
                </Header>
              </div>
            }
            {this.state.resolution.filled && !this.state.resolution.ok &&
              <div>
                <Header className="header" as="h4">
                  Error
                </Header>
                <p>
                  Following error was returned by the server:
                </p>
                <p>
                  <b>Message</b>: {this.state.resolution.error.message.toUpperCase()}
                </p>
                <p>
                  <b>Description</b>: {this.state.resolution.error.description.toUpperCase()}
                </p>
              </div>
            }
        </div>
        );
    }
}

export default Riddle;
