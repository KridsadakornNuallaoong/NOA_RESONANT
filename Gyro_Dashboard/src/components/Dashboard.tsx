import { Checkbox, CheckboxIndicator } from "@radix-ui/react-checkbox";
import React, { useEffect, useState } from "react";
import { FaChartLine, FaRss, FaTemperatureHigh, FaWind } from "react-icons/fa";
import {
    CartesianGrid,
    Legend,
    Line,
    LineChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from "recharts";

const SOCKET_URL = "ws://localhost:5050/ws"; // Replace with your WebSocket server URL

interface AxisData {
    Acceleration: number;
    VelocityAngular: number;
    VibrationSpeed: number;
    VibrationAngle: number;
    VibrationDisplacement: number;
    Frequency: number;
}

interface DataBool {
    Acceleration: boolean;
    VelocityAngular: boolean;
    VibrationSpeed: boolean;
    VibrationAngle: boolean;
    VibrationDisplacement: boolean;
    Frequency: boolean;
    Temperature: boolean;
}

interface GyroData {
    DeviceAddress: string;
    DateTime: string;
    TimeStamp: number;
    X: AxisData;
    Y: AxisData;
    Z: AxisData;
    Temperature: number;
    ModbusHighSpeed: boolean;
}
 
{/* <div className="mt-6">
<h2 className="text-xl font-semibold mb-2 text-center">📈 Acceleration Over Time</h2>
<ResponsiveContainer width="100%" height={300}>
    <LineChart data={data}>
    <CartesianGrid strokeDasharray="3 3" />
    <XAxis dataKey="DateTime" />
    <YAxis />
    <Tooltip />
    <Legend />
    {selectedAxes.X && <Line type="monotone" dataKey="X.Acceleration" stroke="#ff4d4d" strokeWidth={2} />}
    {selectedAxes.Y && <Line type="monotone" dataKey="Y.Acceleration" stroke="#4da6ff" strokeWidth={2} />}
    {selectedAxes.Z && <Line type="monotone" dataKey="Z.Acceleration" stroke="#4dff4d" strokeWidth={2} />}
    </LineChart>
</ResponsiveContainer>
</div> */}

interface GraphProps {
    data: GyroData[];
    label: string;
}

const GraphAxis: React.FC<GraphProps> = ({ data, label }) => {
    const [selectedAxes, setSelectedAxes] = useState({ X: true, Y: true, Z: true });

    const toggleAxis = (axis: "X" | "Y" | "Z") => {
        setSelectedAxes((prev) => ({ ...prev, [axis]: !prev[axis] }));
    };

    return (
        <div className="mt-6">
            {/* Axis Selection */}
            <div className="flex justify-center gap-4 mb-4">
                {["X", "Y", "Z"].map((axis) => (
                <label key={axis} className="flex items-center gap-2 cursor-pointer">
                    <Checkbox
                    checked={selectedAxes[axis as "X" | "Y" | "Z"]}
                    onCheckedChange={() => toggleAxis(axis as "X" | "Y" | "Z")}
                    className="w-5 h-5 border border-[#43c7a2] rounded bg-white checked:bg-blue-500"
                    >
                    <CheckboxIndicator className="w-full h-full flex items-center justify-center">
                        ✅
                    </CheckboxIndicator>
                    </Checkbox>
                    {axis}-Axis
                </label>
                ))}
            </div>
            <h2 className="text-xl font-semibold mb-2 text-center">📈 {label}</h2>
            <ResponsiveContainer width="100%" height={300}>
                <LineChart data={data}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="DateTime" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    {selectedAxes.X && <Line type="monotone" dataKey={"X." + label} stroke="#ff4d4d" strokeWidth={2} />}
                    {selectedAxes.Y && <Line type="monotone" dataKey={"Y." + label} stroke="#4da6ff" strokeWidth={2} />}
                    {selectedAxes.Z && <Line type="monotone" dataKey={"Z." + label} stroke="#4dff4d" strokeWidth={2} />}
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
}

const SingleGraph: React.FC<GraphProps> = ({ data, label }) => {
    return (
        <div className="mt-6">
            <h2 className="text-xl font-semibold mb-2 text-center">📈 Temperature Over Time</h2>
            <ResponsiveContainer width="100%" height={300}>
                <LineChart data={data}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="DateTime" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Line type="monotone" dataKey={label} stroke="#ff4d4d" strokeWidth={2} />
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
}

const Dashboard: React.FC = () => {
  const [data, setData] = useState<GyroData[]>([]);
  const [selectedAxes, setSelectedAxes] = useState({ X: true, Y: true, Z: true });

  const [selectedGraphs, setSelectedGraphs] = useState<DataBool>({
    Acceleration: true,
    VelocityAngular: true,
    VibrationSpeed: true,
    VibrationAngle: true,
    VibrationDisplacement: true,
    Frequency: true,
    Temperature: true
  });

  const handleGraphToggle = (label: string) => {
    setSelectedGraphs((prevState: any) => ({
      ...prevState,
      [label]: !prevState[label]
    }));
  };

  useEffect(() => {
    const websocket = new WebSocket(SOCKET_URL);

    websocket.onopen = () => {
      console.log("WebSocket connected!");
    };

    websocket.onmessage = (event) => {
      try {
        const newData: GyroData = JSON.parse(event.data);
        console.log("Received:", newData);
        setData((prevData) => [...prevData.slice(-200), newData]); // Keep last 100 records
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    };

    websocket.onclose = () => {
      console.warn("WebSocket closed. Reconnecting...");
      setTimeout(() => {
        const newWs = new WebSocket(SOCKET_URL);
        newWs.onmessage = websocket.onmessage;
      }, 3000);
    };

    return () => {
      websocket.close();
    };
  }, []);

  const toggleAxis = (axis: "X" | "Y" | "Z") => {
    setSelectedAxes((prev) => ({ ...prev, [axis]: !prev[axis] }));
  };

  return (
    <div className="p-5 bg-gray-100 min-h-screen flex flex-col">

      {/* Axis Selection */}
      <div className="flex justify-center gap-4 mb-4">
        {["X", "Y", "Z"].map((axis) => (
          <label key={axis} className="flex items-center gap-2 cursor-pointer">
            <Checkbox
              checked={selectedAxes[axis as "X" | "Y" | "Z"]}
              onCheckedChange={() => toggleAxis(axis as "X" | "Y" | "Z")}
              className="w-5 h-5 border border-[#43c7a2] rounded bg-white checked:bg-blue-500"
            >
              <CheckboxIndicator className="w-full h-full flex items-center justify-center">
                ✅
              </CheckboxIndicator>
            </Checkbox>
            {axis}-Axis
          </label>
        ))}
      </div>

      {/* Layout: Graph on left, Table on right */}
      <div className="flex justify-between gap-4">
        {/* Graph Section */}
        <div className="flex-1">
      <div className="flex justify-between mb-6">
        {/* Graph selection icons */}
        <button
          onClick={() => handleGraphToggle("Acceleration")}
          className="flex items-center space-x-2 text-xl"
        >
          <FaChartLine />
          <span>Acceleration</span>
        </button>
        <button
          onClick={() => handleGraphToggle("VelocityAngular")}
          className="flex items-center space-x-2 text-xl"
        >
          <FaWind />
          <span>Angular Velocity</span>
        </button>
        <button
          onClick={() => handleGraphToggle("VibrationSpeed")}
          className="flex items-center space-x-2 text-xl"
        >
          <FaRss />
          <span>Vibration Speed</span>
        </button>
        <button
          onClick={() => handleGraphToggle("Temperature")}
          className="flex items-center space-x-2 text-xl"
        >
          <FaTemperatureHigh />
          <span>Temperature</span>
        </button>
      </div>

      {/* Conditional Rendering of Graphs */}
      {selectedGraphs.Acceleration && (
        <GraphAxis data={data} label="Acceleration" />
      )}
      {selectedGraphs.VelocityAngular && (
        <GraphAxis data={data} label="VelocityAngular" />
      )}
      {selectedGraphs.VibrationSpeed && (
        <GraphAxis data={data} label="VibrationSpeed" />
      )}
      {selectedGraphs.VibrationAngle && (
        <GraphAxis data={data} label="VibrationAngle" />
      )}
      {selectedGraphs.VibrationDisplacement && (
        <GraphAxis data={data} label="VibrationDisplacement" />
      )}
      {selectedGraphs.Frequency && (
        <GraphAxis data={data} label="Frequency" />
      )}
      {selectedGraphs.Temperature && (
        <SingleGraph data={data} label="Temperature" />
      )}
    </div>

        {/* Table Section */}
        <div className="w-1/2 overflow-auto">
          <table className="min-w-full bg-white shadow-md rounded-lg overflow-hidden">
            <thead className="bg-[#43c7a2] text-black">
              <tr>
            <th className="p-2 text-[10px]">Device</th>
            <th className="p-2 text-[10px]">Timestamp</th>
            {selectedAxes.X && <th className="p-2 text-[10px]">Accel (X)</th>}
            {selectedAxes.X && <th className="p-2 text-[10px]">Angular Vel (X)</th>}
            {selectedAxes.X && <th className="p-2 text-[10px]">Vibration Speed (X)</th>}
            {selectedAxes.X && <th className="p-2 text-[10px]">Vibration Angle (X)</th>}
            {selectedAxes.X && <th className="p-2 text-[10px]">Vibration Displacement (X)</th>}
            {selectedAxes.X && <th className="p-2 text-[10px]">Frequency (X)</th>}
            {selectedAxes.Y && <th className="p-2 text-[10px]">Accel (Y)</th>}
            {selectedAxes.Y && <th className="p-2 text-[10px]">Angular Vel (Y)</th>}
            {selectedAxes.Y && <th className="p-2 text-[10px]">Vibration Speed (Y)</th>}
            {selectedAxes.Y && <th className="p-2 text-[10px]">Vibration Angle (Y)</th>}
            {selectedAxes.Y && <th className="p-2 text-[10px]">Vibration Displacement (Y)</th>}
            {selectedAxes.Y && <th className="p-2 text-[10px]">Frequency (Y)</th>}
            {selectedAxes.Z && <th className="p-2 text-[10px]">Accel (Z)</th>}
            {selectedAxes.Z && <th className="p-2 text-[10px]">Angular Vel (Z)</th>}
            {selectedAxes.Z && <th className="p-2 text-[10px]">Vibration Speed (Z)</th>}
            {selectedAxes.Z && <th className="p-2 text-[10px]">Vibration Angle (Z)</th>}
            {selectedAxes.Z && <th className="p-2 text-[10px]">Vibration Displacement (Z)</th>}
            {selectedAxes.Z && <th className="p-2 text-[10px]">Frequency (Z)</th>}
              </tr>
            </thead>
            <tbody>
              {data.slice(-30).map((entry, index) => (
            <tr key={index} className="text-center border-t">
              <td className="p-2">{entry.DeviceAddress}</td>
              <td className="p-2">{entry.DateTime}</td>
              {selectedAxes.X && <td className="p-2">{entry.X.Acceleration.toFixed(2)}</td>}
                {selectedAxes.X && <td className="p-2">{entry.X.VelocityAngular.toFixed(2)}</td>}
                {selectedAxes.X && <td className="p-2">{entry.X.VibrationSpeed.toFixed(2)}</td>}
                {selectedAxes.X && <td className="p-2">{entry.X.VibrationAngle.toFixed(2)}</td>}
                {selectedAxes.X && <td className="p-2">{entry.X.VibrationDisplacement.toFixed(2)}</td>}
                {selectedAxes.X && <td className="p-2">{entry.X.Frequency.toFixed(2)}</td>}
                {selectedAxes.Y && <td className="p-2">{entry.Y.Acceleration.toFixed(2)}</td>}
                {selectedAxes.Y && <td className="p-2">{entry.Y.VelocityAngular.toFixed(2)}</td>}
                {selectedAxes.Y && <td className="p-2">{entry.Y.VibrationSpeed.toFixed(2)}</td>}
                {selectedAxes.Y && <td className="p-2">{entry.Y.VibrationAngle.toFixed(2)}</td>}
                {selectedAxes.Y && <td className="p-2">{entry.Y.VibrationDisplacement.toFixed(2)}</td>}
                {selectedAxes.Y && <td className="p-2">{entry.Y.Frequency.toFixed(2)}</td>}
                {selectedAxes.Z && <td className="p-2">{entry.Z.Acceleration.toFixed(2)}</td>}
                {selectedAxes.Z && <td className="p-2">{entry.Z.VelocityAngular.toFixed(2)}</td>}
                {selectedAxes.Z && <td className="p-2">{entry.Z.VibrationSpeed.toFixed(2)}</td>}
                {selectedAxes.Z && <td className="p-2">{entry.Z.VibrationAngle.toFixed(2)}</td>}
                {selectedAxes.Z && <td className="p-2">{entry.Z.VibrationDisplacement.toFixed(2)}</td>}
                {selectedAxes.Z && <td className="p-2">{entry.Z.Frequency.toFixed(2)}</td>}
            </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
