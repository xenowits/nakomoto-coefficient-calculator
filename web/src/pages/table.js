import React from "react";
import "bootstrap/dist/css/bootstrap.min.css";
import BTable from "react-bootstrap/Table";
import { useTable } from "react-table";

function Table({ columns, data }) {
  // Use the state and functions returned from useTable to build your UI
  const { getTableProps, headerGroups, rows, prepareRow } = useTable({
    columns,
    data,
  });

  const style = {
    maxWidth: 800,
    alignItems: "center",
  };

  // Render the UI for your table
  return (
    <BTable striped bordered hover size="sm" {...getTableProps()} style={style}>
      <thead>
        {headerGroups.map((headerGroup) => (
          <tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map((column) => (
              <th {...column.getHeaderProps()}>{column.render("Header")}</th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody>
        {rows.map((row, i) => {
          prepareRow(row);
          return (
            <tr {...row.getRowProps()}>
              {row.cells.map((cell) => {
                return <td {...cell.getCellProps()}>{cell.render("Cell")}</td>;
              })}
            </tr>
          );
        })}
      </tbody>
    </BTable>
  );
}

function App() {
  const columns = React.useMemo(
    () => [
      {
        Header: "Name",
        columns: [
          {
            Header: "Chain",
            accessor: "chainName",
          },
          {
            Header: "Token",
            accessor: "tokenName",
          },
        ],
      },
      {
        Header: "Coefficient",
        columns: [
          {
            Header: "Current Value",
            accessor: "currVal",
          },
          {
            Header: "Previous Value",
            accessor: "prevVal",
          },
          {
            Header: "Change",
            accessor: "changeVal",
          },
        ],
      },
    ],
    []
  );

  const dataList = [
    {
      chainName: "Binance",
      tokenName: "BNB",
      currVal: 7,
      prevVal: 7,
      changeVal: 0,
    },
    {
      chainName: "Cosmos",
      tokenName: "ATOM",
      currVal: 7,
      prevVal: 7,
      changeVal: 0,
    },
    {
      chainName: "Osmosis",
      tokenName: "OSMO",
      currVal: 4,
      prevVal: 4,
      changeVal: 0,
    },
    {
      chainName: "Polygon",
      tokenName: "MATIC",
      currVal: 2,
      prevVal: 2,
      changeVal: 0,
    },
    {
      chainName: "Mina",
      tokenName: "MINA",
      currVal: 11,
      prevVal: 11,
      changeVal: 0,
    },
  ];

  const data = React.useMemo(() => dataList, []);

  return (
    <div>
      <Table columns={columns} data={data} />
    </div>
  );
}

export default App;
