import React, { ReactElement } from 'react';
import {
    TextInput,
    PageSection,
    Form,
    SelectOption,
    Checkbox,
    TextArea,
} from '@patternfly/react-core';
import * as yup from 'yup';

import FormMultiSelect from 'Components/FormMultiSelect';
import usePageState from 'Containers/Integrations/hooks/usePageState';
import useIntegrationForm from '../useIntegrationForm';
import { IntegrationFormProps } from '../integrationFormTypes';

import IntegrationFormActions from '../IntegrationFormActions';
import FormCancelButton from '../FormCancelButton';
import FormTestButton from '../FormTestButton';
import FormSaveButton from '../FormSaveButton';
import FormMessage from '../FormMessage';
import FormLabelGroup from '../FormLabelGroup';

export type GoogleIntegration = {
    id?: string;
    name: string;
    categories: ('REGISTRY' | 'SCANNER')[];
    google: {
        endpoint: string;
        project: string;
        serviceAccount: string;
    };
    skipTestIntegration: boolean;
    type: 'google';
    enabled: boolean;
    clusterIds: string[];
};

export type GoogleIntegrationFormValues = {
    config: GoogleIntegration;
    updatePassword: boolean;
};

export const validationSchema = yup.object().shape({
    config: yup.object().shape({
        name: yup.string().trim().required('An integration name is required'),
        categories: yup
            .array()
            .of(yup.string().trim().oneOf(['REGISTRY', 'SCANNER']))
            .min(1, 'Must have at least one type selected')
            .required('A category is required'),
        google: yup.object().shape({
            endpoint: yup.string().trim().required('An endpoint is required'),
            project: yup.string().trim().required('A project is required'),
            serviceAccount: yup
                .string()
                .test(
                    'serviceAccount-test',
                    'A service account key is required',
                    (value, context: yup.TestContext) => {
                        const requirePasswordField =
                            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                            // @ts-ignore
                            context?.from[2]?.value?.updatePassword || false;

                        if (!requirePasswordField) {
                            return true;
                        }

                        const trimmedValue = value?.trim();
                        return !!trimmedValue;
                    }
                ),
        }),
        skipTestIntegration: yup.bool(),
        type: yup.string().matches(/google/),
        enabled: yup.bool(),
        clusterIds: yup.array().of(yup.string()),
    }),
    updatePassword: yup.bool(),
});

export const defaultValues: GoogleIntegrationFormValues = {
    config: {
        name: '',
        categories: [],
        google: {
            endpoint: '',
            project: '',
            serviceAccount: '',
        },
        skipTestIntegration: false,
        type: 'google',
        enabled: true,
        clusterIds: [],
    },
    updatePassword: true,
};

function DockerIntegrationForm({
    initialValues = null,
    isEditable = false,
}: IntegrationFormProps<GoogleIntegration>): ReactElement {
    const formInitialValues = defaultValues;
    if (initialValues) {
        formInitialValues.config = { ...formInitialValues.config, ...initialValues };
        // We want to clear the password because backend returns '******' to represent that there
        // are currently stored credentials
        formInitialValues.config.google.serviceAccount = '';
    }
    const {
        values,
        touched,
        errors,
        dirty,
        isValid,
        setFieldValue,
        handleBlur,
        isSubmitting,
        isTesting,
        onSave,
        onTest,
        onCancel,
        message,
    } = useIntegrationForm<GoogleIntegrationFormValues, typeof validationSchema>({
        initialValues: formInitialValues,
        validationSchema,
    });

    const { isCreating } = usePageState();

    function onChange(value, event) {
        return setFieldValue(event.target.id, value);
    }

    function onCustomChange(id, value) {
        return setFieldValue(id, value);
    }

    return (
        <>
            <PageSection variant="light" isFilled hasOverflowScroll>
                {message && <FormMessage message={message} />}
                <Form isWidthLimited>
                    <FormLabelGroup
                        label="Integration name"
                        isRequired
                        fieldId="config.name"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            type="text"
                            id="config.name"
                            placeholder="(ex. Google Registry and Scanner)"
                            value={values.config.name}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Type"
                        isRequired
                        fieldId="config.categories"
                        touched={touched}
                        errors={errors}
                    >
                        <FormMultiSelect
                            id="config.categories"
                            values={values.config.categories}
                            onChange={onCustomChange}
                            isDisabled={!isEditable}
                        >
                            <SelectOption key={0} value="REGISTRY">
                                Registry
                            </SelectOption>
                            <SelectOption key={1} value="SCANNER">
                                Scanner
                            </SelectOption>
                        </FormMultiSelect>
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Registry endpoint"
                        isRequired
                        fieldId="config.google.endpoint"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            type="text"
                            id="config.google.endpoint"
                            placeholder="(ex. gcr.io)"
                            value={values.config.google.endpoint}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Project"
                        isRequired
                        fieldId="config.google.project"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            type="text"
                            id="config.google.project"
                            value={values.config.google.project}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    {!isCreating && (
                        <FormLabelGroup
                            fieldId="updatePassword"
                            helperText="Setting this to false will use the currently stored credentials, if they exist."
                            errors={errors}
                        >
                            <Checkbox
                                id="updatePassword"
                                label="Update stored credentials"
                                isChecked={values.updatePassword}
                                onChange={onChange}
                                onBlur={handleBlur}
                                isDisabled={!isEditable}
                            />
                        </FormLabelGroup>
                    )}
                    <FormLabelGroup
                        isRequired={values.updatePassword}
                        label="Service account key (JSON)"
                        fieldId="config.google.serviceAccount"
                        touched={touched}
                        errors={errors}
                    >
                        <TextArea
                            isRequired={values.updatePassword}
                            id="config.google.serviceAccount"
                            value={values.config.google.serviceAccount}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable || !values.updatePassword}
                            placeholder={
                                values.updatePassword
                                    ? ''
                                    : 'Currently-stored password will be used.'
                            }
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        fieldId="config.skipTestIntegration"
                        touched={touched}
                        errors={errors}
                    >
                        <Checkbox
                            label="Create integration without testing"
                            id="config.skipTestIntegration"
                            isChecked={values.config.skipTestIntegration}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                </Form>
            </PageSection>
            {isEditable && (
                <IntegrationFormActions>
                    <FormSaveButton
                        onSave={onSave}
                        isSubmitting={isSubmitting}
                        isTesting={isTesting}
                        isDisabled={!dirty || !isValid}
                    >
                        Save
                    </FormSaveButton>
                    <FormTestButton
                        onTest={onTest}
                        isSubmitting={isSubmitting}
                        isTesting={isTesting}
                        isValid={isValid}
                    >
                        Test
                    </FormTestButton>
                    <FormCancelButton onCancel={onCancel}>Cancel</FormCancelButton>
                </IntegrationFormActions>
            )}
        </>
    );
}

export default DockerIntegrationForm;
