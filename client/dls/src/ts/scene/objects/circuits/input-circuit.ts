import { SimEvent, UpdationStrategy } from '@ts/engine.js';
import { ConsumerPin } from '../consumer-pin.js';
import { ProducerPin } from '../producer-pin.js';
import { simEngine } from '@src/routes/dls/+page.svelte';
import { type Circuit, circuitCloneHelper, type Props, CircuitPropType } from './circuit.js';
import { CircuitSceneObject } from './circuit.js';

type InputCircuitProps = Props & {
	value: boolean;
};

export class InputCircuit implements Circuit {
	updationStrategy = UpdationStrategy.InNextFrame;
	inputWireUpdationStrategy = UpdationStrategy.InNextFrame;
	outputWireUpdationStrategy = UpdationStrategy.InNextFrame;

	circuitType = 'Input';
	props: InputCircuitProps = {
		label: 'Input',
		value: false
	};

	propTypes = {
		value: CircuitPropType.Bool
	};
	propSetters = {
		value: (circuit: Circuit, value: any) => {
			console.log('setting value prop to: ', value);
			if (typeof value != 'boolean') {
				throw Error();
			}
			circuit.props.value = value;
			circuit.producerPins[0].setValue(circuit.props.value);
			return true;
		}
	};

	simFrameAllocated = false;

	consumerPins: ConsumerPin[];
	producerPins: ProducerPin[];

	sceneObject: CircuitSceneObject | undefined;

	constructor(value: boolean) {
		this.sceneObject = undefined;

		this.consumerPins = new Array();

		this.producerPins = Array(1);
		for (let i = 0; i < this.producerPins.length; i++) {
			this.producerPins[i] = new ProducerPin(this, i);
		}

		this.updateHandeler(this);

		simEngine.recurringEvents.push(new SimEvent(this, this.updateHandeler));
	}

	updateHandeler(self_: Circuit) {
		const self = self_ as InputCircuit;
		self.producerPins[0].setValue(self.props.value);
	}

	clone(): Circuit {
		return circuitCloneHelper(this);
	}

	onSceneObjectConfigured(): void {}
}
